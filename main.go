package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	//Environment variables
	client_id := os.Getenv("SPOTIFY_ID")
	client_secret := os.Getenv("SPOTIFY_SECRET")

	port := ":3000"
	fmt.Println(client_id + "\n" + client_secret + "\n" + "port " + port)
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(port, nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!");
}
