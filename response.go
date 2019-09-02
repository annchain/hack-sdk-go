package main

type QueryNonceResp struct {
	Nonce uint64 `json:"data"`
	Err   string `json:"message"`
}

type QueryBalanceResp struct {
	Data QueryBalanceRespData `json:"data"`
	Err  string               `json:"message"`
}
type QueryBalanceRespData struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type QueryTxsResp struct {
	Data PoolTxs `json:"data"`
	Err  string  `json:"message"`
}
type PoolTxs struct {
	Seq Sequencer     `json:"sequencer"`
	Txs []Transaction `json:"transactions"`
}
