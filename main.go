package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("GoFlix")

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	srv := newServer()
	srv.store = &dbStore{}
	err := srv.store.Open()
	if err != nil {
		return err
	}
	http.HandleFunc("/", srv.serveHTTP)
	defer srv.store.Close()
	fmt.Printf("Server start on port 8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		return err
	}

	return nil
}
