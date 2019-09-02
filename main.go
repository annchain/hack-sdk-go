package main

import "fmt"

func main() {
	url := "http://localhost:8000"

	og, _ := NewOgSolver(url, "")

	// ----------- query pool txs ----------------

	resp, err := og.QueryAllTipsInPool()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("resp:")
	fmt.Println(resp)

	// ----------- query balance ----------------

	//resp, err := og.QueryBalance("0x8b605f016cfe161f66eb7a0d8f97d2a9b098d3cc")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println("resp: ", resp)

	// ----------- gen account ----------------

	//a := GenerateAccount()
	//fmt.Println("priv: ", a.PrivateKey)
	//fmt.Println("pub: ", a.PublicKey)
	//fmt.Println("addr: ", a.Address)
}
