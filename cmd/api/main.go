package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test log")
	fmt.Fprintf(w, "Hello, Mike!")
}
