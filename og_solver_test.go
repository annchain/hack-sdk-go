package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestGenerateAccount(t *testing.T) {
	a := GenerateAccount()

	fmt.Println("priv: ", a.PrivateKey)
	fmt.Println("pub: ", a.PublicKey)
	fmt.Println("addr: ", a.Address)
}

func TestOgSolver_QueryBalance(t *testing.T) {
	url := "http://localhost:8000"

	og, _ := NewOgSolver(url, "")
	resp, err := og.QueryBalance("0x8b605f016cfe161f66eb7a0d8f97d2a9b098d3cc")
	if err != nil {
		fmt.Println(err)
		return
	}

	expectBalance := "1000000"
	if resp != expectBalance {
		t.Fatalf("balance not correct, should be: %s, get: %s", expectBalance, resp)
	}
}

func TestOgSolver_SendTx(t *testing.T) {
	url := "http://localhost:8000"

	priv := "af1b6df8cc06d79902029c0e446c3dc2788893185759d2308b5bb10aa0614b7d"
	og, _ := NewOgSolver(url, priv)

	nonce, err := og.QueryNonce(og.Address())
	if err != nil {
		t.Fatal(err)
	}

	poolTxs, err := og.QueryAllTipsInPool()
	if err != nil {
		t.Fatal(err)
	}
	seq := poolTxs.Seq

	tx := Transaction{
		Parents:   []string{seq.Hash},
		From:      "0xf1b4b3de579ff16888f3340f39c45f207f2cd84d",
		Nonce:     nonce + 1,
		Value:     big.NewInt(0),
		Guarantee: big.NewInt(1),
	}

	msg, _ := tx.SignatureTarget()
	t.Logf("msg: %x", msg)

	hash, err := og.SendTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("hash: %s", hash)
}
