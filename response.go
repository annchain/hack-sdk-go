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

func (tr *TransactionResp) FromMap(m map[string]interface{}) {
	tr.Type = m["type"].(int)
	tr.Hash = m["hash"].(string)
	tr.Parents = m["parents"].([]string)
	tr.From = m["from"].(string)
	tr.To = m["to"].(string)
	tr.Nonce = m["nonce"].(uint64)
	tr.Guarantee = m["guarantee"].(string)
	tr.Value = m["value"].(string)
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

func (tr *SequencerResp) FromMap(m map[string]interface{}) {
	tr.Type = m["type"].(int)
	tr.Hash = m["hash"].(string)
	tr.Parents = m["parents"].([]string)
	tr.From = m["from"].(string)
	tr.Nonce = m["nonce"].(uint64)
	tr.Treasure = m["treasure"].(string)
	tr.Height = m["height"].(uint64)
}

type TxiResp struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}
