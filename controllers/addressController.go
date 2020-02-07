package controllers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
)

type AddressController struct {
	beego.Controller
}

func (this *AddressController) Get() {
	// 警告：本网站未做信息加密处理，请不要将获得到的密钥用于真正的钱包中，以免被盗。
	// Warning: this website does not encrypt information, please do not use the obtained key in the real wallet, in case of theft.


	// Version  Public key hash                           Checksum
	// 00       62E907B15CBF27D5425399EBF6F0FB50EBB88F18  C29B7D93

	// PubkeyHash = RIPEMD160(SHA256(publicKey))
	// checksum = SHA256(SHA256(version + PubkeyHash))

	this.Data["Version"] = "00"
	this.Data["PubKeyHash"] = "9b96f76f329936e02dde9cb7ed2a93fcbd4f9809"
	this.Data["Checksum"] = "82800b37"

	this.Data["Address"] = "1FBgb5E2jfcWMKbLjhdBuMfm3nGYpp24qt"

	this.Data["PublicKey"] = "9829588ae2efbd5eaa879b4db7703b85b257f4d2594b641cef9d5cebd29ce36799f49a195d660df5f12fc9131af4c5607ace2161cf4639e3dd3559b1f91cf7bf"
	this.Data["PrivateKey"] = "5b5c060bd36f85e6587f00a16e697f2bf8bf0ac08e3c40d8c78275b2122b1c8d"

	this.TplName = "address.html"
}

func (this *AddressController) Post() {
	// form data

	wallete := NewWallet()

	address := wallete.GetAddress()

	pubKeyHash := Base58Decode(address)
	version := pubKeyHash[:1]
	checksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash) - addressChecksumLen]

	this.Data["Version"] = fmt.Sprintf("%x", version)
	this.Data["PubKeyHash"] = fmt.Sprintf("%x", pubKeyHash)
	this.Data["Checksum"] = fmt.Sprintf("%x", checksum)

	this.Data["Address"] = fmt.Sprintf("%s", address)
	this.Data["PublicKey"] = fmt.Sprintf("%x", wallete.PublicKey)
	this.Data["PrivateKey"] = fmt.Sprintf("%x", wallete.PrivateKey.D.Bytes())

	this.TplName = "address.html"
}

const (
	version = byte(0x00)
	addressChecksumLen = 4
)

// Wallet 存储公钥和私钥
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet 创建并返回一个钱包
func NewWallet() *Wallet {
	private, pubKey := newKeyPair()

	return &Wallet{
		PrivateKey: private,
		PublicKey:  pubKey,
	}
}

// GetAddress 返回钱包地址

// e.g.
// Version  Public key hash                           Checksum
// 00       62E907B15CBF27D5425399EBF6F0FB50EBB88F18  C29B7D93
func (w *Wallet) GetAddress() []byte {
	// Address = base58(version + pubKeyHash + checksum)
	pubKeyHash := HashPubKey(w.PublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionPayload)

	fullPayload := append(versionPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// HashPubKey hashes public key
// HashPubKey 返回公钥的哈希
func HashPubKey(pubKey []byte) []byte {
	// RIPEMD160(SHA256(PubKey))
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// 生成一个密钥对
// 基于椭圆曲线的算法中, 公钥是曲线上的点, 因此公钥是X, Y坐标的组合
func newKeyPair() (privateKey ecdsa.PrivateKey, publicKey []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// Checksum 根据公钥生成校验和
func checksum(payload []byte) []byte {
	// SHA256(SHA256(payload))
	FirstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(FirstSHA[:])

	// 只取前四个字节
	return secondSHA[:addressChecksumLen]
}

// ValidateAddress 检测地址是否有效
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	if len(pubKeyHash) < addressChecksumLen {
		return false
	}
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	targetChecksum := checksum(pubKeyHash[:len(pubKeyHash)-addressChecksumLen])

	return bytes.Compare(targetChecksum, actualChecksum) == 0
}


// Base58
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	// https://en.bitcoin.it/wiki/Base58Check_encoding#Version_bytes
	if input[0] == 0x00 {
		result = append(result, b58Alphabet[0])
	}

	ReverseBytes(result)

	return result
}

// Base58Decode decodes Base58-encoded data
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)

	for _, b := range input {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()

	if input[0] == b58Alphabet[0] {
		decoded = append([]byte{0x00}, decoded...)
	}

	return decoded
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
