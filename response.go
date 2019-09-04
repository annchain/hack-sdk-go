package main

type NewTxResp struct {
	Hash string `json:"data"`
	Err  string `json:"err"`
}

type QueryNonceResp struct {
	Nonce uint64 `json:"data"`
	Err   string `json:"err"`
}

type QueryBalanceResp struct {
	Data QueryBalanceRespData `json:"data"`
	Err  string               `json:"err"`
}
type QueryBalanceRespData struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type QueryPoolTxsResp struct {
	Data PoolTxs `json:"data"`
	Err  string  `json:"err"`
}
type PoolTxs struct {
	Seq SequencerResp     `json:"sequencer"`
	Txs []TransactionResp `json:"transactions"`
}

type QueryTxsResp struct {
	Data []TransactionResp `json:"data"`
	Err  string            `json:"err"`
}

type QueryTransactionResp struct {
	Data TransactionResp `json:"data"`
	Err  string          `json:"err"`
}

type TransactionResp struct {
	Type      int      `json:"type"`
	Hash      string   `json:"hash"`
	Parents   []string `json:"parents"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Nonce     uint64   `json:"nonce"`
	Guarantee string   `json:"guarantee"`
	Value     string   `json:"value"`
}

type QuerySequencerResp struct {
	Data SequencerResp `json:"data"`
	Err  string        `json:"err"`
}

type SequencerResp struct {
	Type     int      `json:"type"`
	Hash     string   `json:"hash"`
	Parents  []string `json:"parents"`
	From     string   `json:"from"`
	Nonce    uint64   `json:"nonce"`
	Treasure string   `json:"treasure"`
	Height   uint64   `json:"height"`
}
