package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OgSolver struct {
	url     string
	account *OgAccount
}

func NewOgSolver(url string, privHex string) (*OgSolver, error) {
	og := &OgSolver{}
	og.url = url

	acc, err := NewAccount(privHex)
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

func (o *OgSolver) SendTx() {

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

	var poolResp QueryTxsResp
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
	resp, err := http.Get(url)
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

func (o *OgSolver) doPostRequest(url string, reqBody map[string]string) ([]byte, error) {
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body error: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyData))
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
