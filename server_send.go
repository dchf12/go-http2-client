package main

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"
)

var html []byte

// HTMLをブラウザに送信
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write(html)
}

// 素数をブラウザに送信
func handlerPrimeSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	//接続断の検知用にコンテキストを取得
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var num int64 = 1
	for id := 1; id <= 100; id++ {
		// 通信が切れても終了
		select {
		case <-ctx.Done():
			fmt.Println("Connection closed from client")
			return
		default:
			// do nothing
		}
		for {
			num++
			//素数を求める
			if big.NewInt(num).ProbablyPrime(20) {
				fmt.Println(num)
				fmt.Fprintf(w, "data: {\"id\": %d, \"number\": %d}\n\n", id, num)
				flusher.Flush()
				time.Sleep(10 * time.Millisecond)
				break
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Connection closed from server")
}
func main() {
	var err error
	html, err = os.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/prime", handlerPrimeSSE)
	fmt.Println("start http listening :18888")
	err = http.ListenAndServe("localhost:18888", nil)
	fmt.Println(err)
}
