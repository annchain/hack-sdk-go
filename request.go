package main

type TransactionReq struct {
	Parents   []string `json:"parents"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Nonce     uint64   `json:"nonce"`
	Guarantee string   `json:"guarantee"`
	Value     string   `json:"value"`
	Signature string   `json:"signature"`
	PublicKey string   `json:"pubkey"`
}

func NewTransactionReq(rawTx Transaction, sig string, publicKey string) TransactionReq {
	return TransactionReq{
		Parents:   rawTx.Parents,
		From:      rawTx.From,
		To:        rawTx.To,
		Nonce:     rawTx.Nonce,
		Guarantee: rawTx.Guarantee.String(),
		Value:     rawTx.Value.String(),
		Signature: sig,
		PublicKey: publicKey,
	}
}
