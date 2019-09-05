package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"time"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	//url := "http://localhost:8000"
	//
	//og, _ := NewOgSolver(url, "")

	// connect
	ws, err := websocket.Dial("ws://localhost:19000", "", "http://localhost/")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// receive
	var m Message
	for {
		time.Sleep(time.Second * 1)

		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			fmt.Println("Error receiving message: ", err.Error())
			continue
		}
		fmt.Println("Message: ", m)
	}

	//origin := "http://localhost/"
	//url := "ws://localhost:19000"
	//ws, err := websocket.Dial(url, "", origin)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
	////	log.Fatal(err)
	////}
	//
	//fmt.Println("sleep 6 seconds")
	//time.Sleep(time.Second * 6)
	//
	//var msg = make([]byte, 512)
	//var n int
	//if n, err = ws.Read(msg); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Received: %s.\n", msg[:n])

}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}
