package hackSDK

import (
	"encoding/json"
	"fmt"
)

func main() {
	url := "http://localhost:8000"
	kafkaUrl := "localhost:9092"
	priv := "af1b6df8cc06d79902029c0e446c3dc2788893185759d2308b5bb10aa0614b7d"

	og, _ := NewOgSolver(url, kafkaUrl, priv)
	c := og.ReceiveNewestTx()

	for {
		fmt.Println("start consuming one data")
		select {
		case txiResp := <-c:
			var jsonData string
			if txiResp.Type == TxBaseTypeNormal {
				txResp := txiResp.Data.(TransactionResp)
				data, _ := json.Marshal(txResp)
				jsonData = string(data)
			}
			if txiResp.Type == TxBaseTypeSequencer {
				seqResp := txiResp.Data.(SequencerResp)
				data, _ := json.Marshal(seqResp)
				jsonData = string(data)
			}

			fmt.Println(jsonData)
		}
	}

}
