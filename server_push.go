package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var image []byte

func init() {
	var err error
	image, err = ioutil.ReadFile("./image.png")
	if err != nil {
		panic(err)
	}
}

func handlerHtml(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		if err := pusher.Push("/image", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

func handlerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(image)
}

func main() {
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/image", handlerImage)
	fmt.Println("start http listening :18443")
	err := http.ListenAndServeTLS("localhost:18443", "server.crt", "server.key", nil)
	if err != nil {
		fmt.Println(err)
	}
}
