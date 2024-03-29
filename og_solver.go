package hackSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

const (
	TxTypeNormal int = iota
	TxTypeSequencer
)

type OgSolver struct {
	url        string
	kafkaUrl   string
	kafkaTopic string
	token      string
	account    *OgAccount
}

func NewOgSolver(url, kafkaUrl, privHex, token string) (*OgSolver, error) {
	og := &OgSolver{}
	og.url = url
	og.kafkaUrl = kafkaUrl
	og.kafkaTopic = "hack-final-test"
	og.token = token

	acc, err := newAccount(privHex)
	if err != nil {
		return nil, err
	}
	og.account = acc

	return og, nil
}

func (o *OgSolver) PrivateKey() string {
	return o.account.PrivateKey
}

func (o *OgSolver) PublicKey() string {
	return o.account.PublicKey
}

func (o *OgSolver) Address() string {
	return o.account.Address
}

func (o *OgSolver) SendTx(tx Transaction) (string, error) {
	url := o.url + "/new_transaction"

	priv, _ := HexToBytes(o.PrivateKey())
	sig, err := tx.Sign(priv)
	if err != nil {
		return "", fmt.Errorf("sign error: %v", err)
	}
	sigStr := BytesToHex(sig)

	txReq := NewTransactionReq(tx, sigStr, o.PublicKey())
	resp, err := o.doPostRequest(url, txReq)
	if err != nil {
		return "", fmt.Errorf("send tx error: %v", err)
	}

	var hashResp NewTxResp
	err = json.Unmarshal(resp, &hashResp)
	if err != nil {
		return "", fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if hashResp.Err != "" {
		return "", fmt.Errorf("server error: %s", hashResp.Err)
	}

	return hashResp.Hash, nil
}

func (o *OgSolver) ReceiveNewestTx() <-chan *TxiResp {
	c := make(chan *TxiResp)

	go o.kafkaConsume(c, sarama.OffsetNewest)
	return c
}

func (o *OgSolver) ReceiveOldestTx() <-chan *TxiResp {
	c := make(chan *TxiResp)

	go o.kafkaConsume(c, sarama.OffsetOldest)
	return c
}

func (o *OgSolver) QueryNonce(address string) (uint64, error) {
	url := o.url + "/query_nonce?address=" + address

	resp, err := o.doGetRequest(url)
	if err != nil {
		return 0, fmt.Errorf("get nonce error: %v", err)
	}

	var nonceResp QueryNonceResp
	err = json.Unmarshal(resp, &nonceResp)
	if err != nil {
		return 0, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if nonceResp.Err != "" {
		return 0, fmt.Errorf("server error: %s", nonceResp.Err)
	}

	return nonceResp.Nonce, nil
}

func (o *OgSolver) QueryBalance(address string) (string, error) {
	url := o.url + "/query_balance?address=" + address

	resp, err := o.doGetRequest(url)
	if err != nil {
		return "", fmt.Errorf("get balance error: %v", err)
	}

	var balanceResp QueryBalanceResp
	err = json.Unmarshal(resp, &balanceResp)
	if err != nil {
		return "", fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if balanceResp.Err != "" {
		return "", fmt.Errorf("server error: %s", balanceResp.Err)
	}

	return balanceResp.Data.Balance, nil
}

func (o *OgSolver) QueryTransaction(hash string) (*TransactionResp, error) {
	url := o.url + "/transaction?hash=" + hash

	resp, err := o.doGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("get transaction error: %v", err)
	}

	var txResp QueryTransactionResp
	err = json.Unmarshal(resp, &txResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if txResp.Err != "" {
		return nil, fmt.Errorf("server error: %s", txResp.Err)
	}

	return &txResp.Data, nil
}

func (o *OgSolver) QuerySequencerByHash(hash string) (*SequencerResp, error) {
	url := o.url + "/sequencer?hash=" + hash
	return o.querySequencer(url)
}

func (o *OgSolver) QuerySequencerByHeight(height uint64) (*SequencerResp, error) {
	url := o.url + "/sequencer?height=" + strconv.Itoa(int(height))
	return o.querySequencer(url)
}

func (o *OgSolver) querySequencer(url string) (*SequencerResp, error) {
	resp, err := o.doGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("get sequencer error: %v", err)
	}

	var seqResp QuerySequencerResp
	err = json.Unmarshal(resp, &seqResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if seqResp.Err != "" {
		return nil, fmt.Errorf("server error: %s", seqResp.Err)
	}

	return &seqResp.Data, nil
}

func (o *OgSolver) QueryNextSequencerInfo() (*QueryNextSeqRespData, error) {
	url := o.url + "/query_next_seq"
	resp, err := o.doGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("get next seq info error: %v", err)
	}

	var nextSeqResp QueryNextSeqResp
	err = json.Unmarshal(resp, &nextSeqResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if nextSeqResp.Err != "" {
		return nil, fmt.Errorf("server error: %s", nextSeqResp.Err)
	}

	return &nextSeqResp.Data, nil

}

func (o *OgSolver) QueryTxsByAddress(address string) ([]TransactionResp, error) {
	url := o.url + "/transactions?address=" + address
	return o.queryTxs(url)
}

func (o *OgSolver) QueryTxsByHeight(height uint64) ([]TransactionResp, error) {
	url := o.url + "/transactions?height=" + strconv.Itoa(int(height))
	return o.queryTxs(url)
}

func (o *OgSolver) queryTxs(url string) ([]TransactionResp, error) {
	resp, err := o.doGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("get txs error: %v", err)
	}

	var txsResp QueryTxsResp
	err = json.Unmarshal(resp, &txsResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if txsResp.Err != "" {
		return nil, fmt.Errorf("server error: %s", txsResp.Err)
	}

	return txsResp.Data, nil
}

func (o *OgSolver) QueryTxNumByHeight(height uint64) (int, error) {
	url := o.url + "/query_tx_num?height=" + strconv.Itoa(int(height))

	resp, err := o.doGetRequest(url)
	if err != nil {
		return 0, fmt.Errorf("get sequencer error: %v", err)
	}

	var txNumResp TxNumResp
	err = json.Unmarshal(resp, &txNumResp)
	if err != nil {
		return 0, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if txNumResp.Err != "" {
		return 0, fmt.Errorf("server error: %s", txNumResp.Err)
	}

	return txNumResp.Data, nil
}

func (o *OgSolver) QueryAllTipsInPool() (*PoolTxs, error) {
	url := o.url + "/query_pool_tips"
	return o.queryPoolTxs(url)
}

func (o *OgSolver) QueryAllTxsInPool() (*PoolTxs, error) {
	url := o.url + "/query_pool_txs"
	return o.queryPoolTxs(url)
}

func (o *OgSolver) queryPoolTxs(url string) (*PoolTxs, error) {
	resp, err := o.doGetRequest(url)
	if err != nil {
		return nil, fmt.Errorf("get txs error: %v", err)
	}

	var poolResp QueryPoolTxsResp
	err = json.Unmarshal(resp, &poolResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response to json error: %v", err)
	}
	if poolResp.Err != "" {
		return nil, fmt.Errorf("server error: %s", poolResp.Err)
	}

	return &poolResp.Data, nil
}

func (o *OgSolver) doGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create GET request error: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: "token", Value: o.token})

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do GET request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}

	return body, nil
}

func (o *OgSolver) doPostRequest(url string, reqBody interface{}) ([]byte, error) {
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyData))
	if err != nil {
		return nil, fmt.Errorf("create GET request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: o.token})

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do POST request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}

	return body, nil
}

func (o *OgSolver) kafkaConsume(receiver chan *TxiResp, offset int64) {
	consumer, err := sarama.NewConsumer([]string{o.kafkaUrl}, nil)
	if err != nil {
		log.Printf("create consumer error: %v\n", err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(o.kafkaTopic, 0, offset)
	if err != nil {
		log.Printf("create partition consumer error: %v\n", err)
		return
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)

			value := msg.Value

			var txiResp TxiResp
			json.Unmarshal(value, &txiResp)

			if txiResp.Type == TxTypeNormal {
				txMap := txiResp.Data.(map[string]interface{})
				tx := TransactionResp{}
				tx.fromMap(txMap)

				txiResp.Data = tx
				receiver <- &txiResp
				continue
			}
			if txiResp.Type == TxTypeSequencer {
				seqMap := txiResp.Data.(map[string]interface{})
				seq := SequencerResp{}
				seq.fromMap(seqMap)

				txiResp.Data = seq
				receiver <- &txiResp
				continue
			}

		}
	}
}
