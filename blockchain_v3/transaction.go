package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TXID) == 0 && tx.Vin[0].Vout == -1
}

// todo 待测试
func (tx Transaction) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func DeserializeTransation(data []byte) Transaction {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var tx Transaction
	err := decoder.Decode(tx)
	if err != nil {
		log.Panic(err)
	}

	return tx
}

func (tx Transaction) Hash() []byte {
	var hash [32]byte
	tx.ID = []byte{}

	hash = sha256.Sum256(tx.Serialize())

	return hash[:]
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TXID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TXID)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	// 以下做法被证实是错误的
	// vin.PubKey 与 vin.Signature 并未清空。 := range 得到的value是值类型
	//txCopy1 := *tx
	//for _, vin := range txCopy1.Vin {
	//	vin.PubKey = nil
	//	vin.Signature = nil
	//}

	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.TXID, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TXID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	curve := elliptic.P256()
	txCopy := tx.TrimmedCopy()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TXID)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:sigLen / 2])
		s.SetBytes(vin.Signature[sigLen / 2:])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:keyLen / 2])
		y.SetBytes(vin.PubKey[keyLen / 2:])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			fmt.Printf("ERROR: Transaction signature is not correct\n")
			return false
		}

		txCopy.Vin[inID].PubKey = nil
	}

	return true
}

func NewCoinbaseTX(to, data string) *Transaction {
	if ValidateAddress(to) == false {
		log.Panic("Address is not correct")
	}
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{TXID: []byte{}, Vout: -1, Signature: nil, PubKey: []byte(data)}
	txout := NewTXOutput(subsidy, to)

	tx := Transaction{ID:nil, Vin: []TXInput{txin}, Vout: []TXOutput{*txout}}
	tx.ID = tx.Hash()

	return &tx
}

// NewUTXOTransaction 创建一个新的交易
func NewUTXOTransaction(wallet *Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs  []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough money")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	from := fmt.Sprintf("%s", wallet.GetAddress())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc - amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

// String 返回一个人类可读的代表交易的字符串
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       ID:      %x", input.TXID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}