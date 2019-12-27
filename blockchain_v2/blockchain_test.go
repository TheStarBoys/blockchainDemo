package blockchain_v2

import (
	"testing"
)

func TestBlockchain(t *testing.T) {
	bc := NewBlockChain("Alice")
	if balance := bc.GetBalance("Alice"); balance != reward {
		t.Errorf("Alice Balance get %f, expected: %f", balance, reward)
	}
}
