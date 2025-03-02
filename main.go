package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

/*
Handler function for /
*/
func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// inspecting query string of requests
	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	// reading request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s, body:\n%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second,
		body)
	io.WriteString(w, "This is my website!\n")
}

/*
Handler function for /hello
*/
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello HTTP!\n")
}

func runServer(s *http.Server, sName string, canceler context.CancelFunc) {
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("%s closed\n", sName)
	} else if err != nil {
		fmt.Printf("error listening for %s: %s\n", err, sName)
	}
	canceler()
}

func main() {
	// make own custom HTTP request mux(multiplexer)
	mux := http.NewServeMux()
	// register handler functions for given patterns
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	// setting context
	ctx, cancelCtx := context.WithCancel(context.Background())

	// setting server 1
	serverOne := &http.Server{
		Addr:    "127.0.0.1:3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	// setting server 2
	serverTwo := &http.Server{
		Addr:    "127.0.0.1:4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go runServer(serverOne, "server one", cancelCtx)
	go runServer(serverTwo, "server two", cancelCtx)

	<-ctx.Done()
}
