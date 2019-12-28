package main

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	testAddress1 = "1FWj3yPEZvEfdGDzamiHAfagksFx2KU3Hu"
	testAddress2 = "19attaBUc7XgVRiZqk4hT2aWfeatDuBQAQ"
)
func TestWallet(t *testing.T) {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("address: %s\n", address)
	if ValidateAddress(string(address)) == false {
		t.Errorf("Address invalid, expected: valid")
	}

	if ValidateAddress("xsfdloOL") == true {
		t.Errorf("Address valid, expected: invalid")
	}

	pubKeyHash := HashPubKey(wallet.PublicKey)
	if bytes.Compare(pubKeyHash, HashPubKey(wallet.PublicKey)) != 0 {
		t.Errorf("HashPubKey ERROR, second hash not equal")
	}

	sum := checksum([]byte("0062E907B15CBF27D5425399EBF6F0FB50EBB88F18"))
	if bytes.Compare(sum, checksum([]byte("0062E907B15CBF27D5425399EBF6F0FB50EBB88F18"))) != 0 {
		t.Errorf("checksum ERROR, second hash not equal")
	}
}
