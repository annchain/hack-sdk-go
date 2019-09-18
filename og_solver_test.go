package hackSDK

import (
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestGenerateAccount(t *testing.T) {
	//a := GenerateAccount()
	//
	//fmt.Println("priv: ", a.PrivateKey)
	//fmt.Println("pub: ", a.PublicKey)
	//fmt.Println("addr: ", a.Address)

	url := "http://localhost:8000"
	kafkaUrl := "localhost:9092"
	priv := "af1b6df8cc06d79902029c0e446c3dc2788893185759d2308b5bb10aa0614b7d"
	og, _ := NewOgSolver(url, kafkaUrl, priv, "")

	fmt.Println(og.Address())
}

func TestOgSolver_QueryBalance(t *testing.T) {
	url := "http://localhost:8000"
	kafkaUrl := "localhost:9092"

	og, _ := NewOgSolver(url, kafkaUrl, "", "")
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

func TestOgSolver_QueryTransaction(t *testing.T) {
	url := "http://localhost:8000"
	kafkaUrl := "localhost:9092"
	priv := "af1b6df8cc06d79902029c0e446c3dc2788893185759d2308b5bb10aa0614b7d"
	og, _ := NewOgSolver(url, kafkaUrl, priv, "")

	hash := "0xa80f781e993539ca0b9b76696a1aab3e5b39e3290cdc85840ae3b90694a25e55"
	tx, err := og.QueryTransaction(hash)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("tx: %v", tx)
}

func TestOgSolver_SendTx(t *testing.T) {
	url := "http://47.100.122.212:30020"
	kafkaUrl := "47.100.222.11:30000"

	priv := "af1b6df8cc06d79902029c0e446c3dc2788893185759d2308b5bb10aa0614b7d"
	token := "98765467890"

	og, _ := NewOgSolver(url, kafkaUrl, priv, token)

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

func TestOgSolver_KafkaConsumer(t *testing.T) {
	url := "http://47.100.122.212:30020"
	kafkaUrl := "47.100.222.11:30000"

	priv := "af1b6df8cc06d79902029c0e446c3dc2788893185759d2308b5bb10aa0614b7d"
	token := "98765467890"

	og, _ := NewOgSolver(url, kafkaUrl, priv, token)

	c := og.ReceiveNewestTx()

	timer := time.NewTimer(time.Second * 15)

	for {
		select {
		case txi := <-c:
			fmt.Println(txi)
		case <-timer.C:
			return
		}
	}
}

func TestOgSolver_TxTarget(t *testing.T) {
	tx := Transaction{}
	tx.From = "c4321fee1e29b13b042feab06dea55e7caf85948"
	tx.To = EmptyAddress
	tx.Nonce = 115
	tx.Parents = []string{"0xba001c4e6b416268761c068a711266e05dd994e804952b0f1033025bfbb54703"}
	tx.Guarantee = big.NewInt(200)
	tx.Value = big.NewInt(0)

	msg, err := tx.SignatureTarget()
	if err != nil {
		t.Fatalf("get signature target error: %v", err)
	}
	t.Logf("msg: %x", msg)

	t.Logf("bigint bytes: %x", big.NewInt(200).Bytes())
}
