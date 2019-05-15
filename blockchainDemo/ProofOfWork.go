package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)
// 前导零(二进制)
const targetBits = 8
type ProofOfWork struct {
	block 	*Block
	target	*big.Int		// 目标值, 因为hash是256位
}
func NewPoW(block *Block) *ProofOfWork{
	// 0x00000000....0001		256位
	target := big.NewInt(1)
	// 左移 256位 - targetBits
	// 前导零 = targetBits / 4
	target.Lsh(target, uint(256 - targetBits))
	pow := ProofOfWork{block:block, target:target}
	return &pow
}

// 挖矿主程序
func (pow *ProofOfWork)Mining() {
	var hashInt big.Int
	var nonce int64
	block := pow.block
	fmt.Println("Begin Mining...")
	fmt.Printf("target hash: %x\n", pow.target.Bytes())
	// 增加nonce值, 直到使hash小于目标值
	for nonce < math.MaxInt64 {
		// 准备数据 nonce改变, data随之改变
		data := pow.Prepare(nonce)
		// func Sum256(data []byte) [Size]byte {
		// 计算hash值
		hash := sha256.Sum256(data)
		// 将hash转成big.Int并返回给自身
		hashInt.SetBytes(hash[:])
		// 如果hash值小于目标值 返回 -1
		if hashInt.Cmp(pow.target) == -1{
			block.Nonce = nonce
			block.Hash = hash[:]
			return
		}
		nonce++
	}

}
// nonce是传进去的
func (pow *ProofOfWork)Prepare(nonce int64) (data []byte){
	var tmp [][]byte
	block := pow.block
	// 不应该包含Data 和 hash, Data数据通过MerkleRoot保证无法篡改, Hash是要求的参数
	tmp = [][]byte{
		Int2Bytes(block.Version),
		Int2Bytes(block.TimeStamp),
		block.MerkleRoot,
		block.PrevBlockHash,
		Int2Bytes(block.Bits),
		Int2Bytes(nonce),
	}
	// 将数据[][]byte 用空[]byte拼接成[]byte
	data = bytes.Join(tmp, []byte{})
	return data
}