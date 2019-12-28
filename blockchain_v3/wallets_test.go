package main

import (
	"testing"
)

func TestWallets(t *testing.T) {
	wallets, err := NewWallets("3000")
	if err != nil && len(wallets.Wallets) != 0 {
		t.Errorf("Wallets File Not Found, but wallets is not empty")
	}
	address := wallets.CreateWallet()
	wallet := wallets.GetWallet(address)
	if address != string(wallet.GetAddress()) {
		t.Errorf("wallet address not equal, address: %s, wallet address: %s", address, wallet.GetAddress())
	}
	addrs := wallets.GetAddresses()
	hasAddr := func() bool {
		for _, addr := range addrs {
			if addr == address {
				return true
			}
		}

		return false
	}()

	if !hasAddr {
		t.Errorf("wallets have not address: %s", address)
	}

	wallets.SaveToFile("3000")

	afterWallets := &Wallets{}
	afterWallets.LoadFromFile("3000")
	afterAddrs := afterWallets.GetAddresses()

	if len(afterAddrs) != len(addrs) {
		t.Errorf("after address length not equal addrs")
	}

	for i := 0; i < len(addrs); i ++ {
		if afterAddrs[i] != addrs[i] {
			t.Errorf("address list not equal")
		}
	}
}
