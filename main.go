package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
	Handler function for /
*/
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

/*
	Handler function for /hello
*/
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello HTTP!\n")
}

func main() {
	// make own custom HTTP request mux(multiplexer)
    mux := http.NewServeMux()
	// register handler functions for given patterns
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	err := http.ListenAndServe("127.0.0.1:3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
