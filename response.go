package hackSDK

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

func (t *TransactionResp) FromMap(m map[string]interface{}) {
	t.Type = int(m["type"].(float64))
	t.Hash = m["hash"].(string)

	var parents []string
	parentsInterface := m["parents"].([]interface{})
	for _, p := range parentsInterface {
		parents = append(parents, p.(string))
	}
	t.Parents = parents

	t.From = m["from"].(string)
	t.To = m["to"].(string)
	t.Nonce = uint64(m["nonce"].(float64))
	t.Guarantee = m["guarantee"].(string)
	t.Value = m["value"].(string)
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

func (s *SequencerResp) FromMap(m map[string]interface{}) {
	s.Type = int(m["type"].(float64))
	s.Hash = m["hash"].(string)

	var parents []string
	parentsInterface := m["parents"].([]interface{})
	for _, p := range parentsInterface {
		parents = append(parents, p.(string))
	}
	s.Parents = parents

	s.From = m["from"].(string)
	s.Nonce = uint64(m["nonce"].(float64))
	s.Treasure = m["treasure"].(string)
	s.Height = uint64(m["height"].(float64))
}

type TxiResp struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}
