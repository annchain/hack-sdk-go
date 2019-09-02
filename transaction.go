package main

type Transaction struct {
	Type      string   `json:"type"`
	Hash      string   `json:"hash"`
	Parents   []string `json:"parents"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Nonce     uint64   `json:"nonce"`
	Guarantee string   `json:"guarantee"`
	Value     string   `json:"value"`
}

type Sequencer struct {
	Type     string   `json:"type"`
	Hash     string   `json:"hash"`
	Parents  []string `json:"parents"`
	From     string   `json:"from"`
	Nonce    uint64   `json:"nonce"`
	Treasure string   `json:"treasure"`
	Height   uint64   `json:"height"`
}
