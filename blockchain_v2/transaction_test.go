package blockchain_v2

import "testing"

func TestTransaction(t *testing.T) {
	bc := NewBlockChain("Alice")
	balance := bc.GetBalance("Alice")
	if balance != reward {
		t.Errorf("Alice Balance get %f, expected: %f", balance, reward)
	}
	coinbase := NewCoinbaseTx("Alice", genesisInfo)
	if coinbase.IsCoinbase() == false {
		t.Errorf("It's not coinbase, expected: coinbase")
	}

	tx := NewTransaction("Alice", "Bob", 5, bc)
	if tx == nil {
		t.Errorf("Alice not have enough money, expected enough money")
	}
	bc.AddBlock([]*Transaction{coinbase, tx})
	balance = bc.GetBalance("Bob")
	if balance != 5 {
		t.Errorf("Bob have %f BTC, expected: 5", balance)
	}
	balance = bc.GetBalance("Alice")
	if balance != (12.5 * 2) - 5 {
		t.Errorf("Alice have %f BTC, expected: %d", balance, 20)
	}

	if bc.Send("Alice", "Failed data", "Ivan", "Bob", 5) == true {
		t.Errorf("Send Error, Ivan not have enough money, expected: %v", false)
	}
}
