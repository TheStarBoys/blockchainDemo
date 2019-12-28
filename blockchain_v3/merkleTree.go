package main

import (
	"crypto/sha256"
	"fmt"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, node)
	}

	for i := 0; i < len(data)/2; i++ {
		var newLevel []*MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}

		nodes = newLevel
	}

	return &MerkleTree{nodes[0]}
}


func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := &MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return mNode
}

func (tree *MerkleTree) ToShow() []string {
	var result []string

	queue := make([]*MerkleNode, 0)
	queue = append(queue, tree.RootNode)

	for len(queue) > 0 {
		length := len(queue)
		for i := 0; i < length; i++ {
			node := queue[0]
			queue = queue[1:]

			result = append(result, fmt.Sprintf("%x", node.Data))

			if node.Left != nil {
				queue = append(queue, node.Left)
			}

			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}

	return result
}

func printTree(tree *MerkleTree) {
	queue := make([]*MerkleNode, 0)
	queue = append(queue, tree.RootNode)

	level := 1
	for len(queue) > 0 {
		length := len(queue)
		fmt.Printf("---------Level%d---------\n", level)
		for i := 0; i < length; i++ {
			node := queue[0]
			queue = queue[1:]
			fmt.Printf("%x ", node.Data)

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		level++
		fmt.Println()
	}
}

func (tree *MerkleTree) PrintTree() {
	printTree(tree)
}