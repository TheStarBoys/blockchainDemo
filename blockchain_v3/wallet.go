package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

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