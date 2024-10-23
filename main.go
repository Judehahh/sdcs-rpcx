package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"rpc"
)

func router(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	xclient := rpc.NewXClient(10086)

	switch r.Method {
	case http.MethodGet:
		key := rpc.Key{Key: r.URL.Path[1:]}
		value := &rpc.Value{}
		if err := xclient.Call(context.Background(), "Get", key, value); err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		if value.Value == "" {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value.Value))
		}
	case http.MethodPost:
		content := string(body)

		key_start := strings.Index(content, `"`)
		if key_start == -1 {
			http.Error(w, "No opening quote found", http.StatusBadRequest)
			return
		}
		key_end := strings.Index(content[key_start+1:], `"`)
		if key_end == -1 {
			http.Error(w, "No closing quote found", http.StatusBadRequest)
			return
		}

		pair := rpc.Pair{
			Key:   content[key_start+1 : key_start+1+key_end],
			Value: content,
		}
		flag := &rpc.Flag{}
		if err := xclient.Call(context.Background(), "Post", pair, flag); err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		key := rpc.Key{Key: r.URL.Path[1:]}
		flag := &rpc.Flag{}
		if err := xclient.Call(context.Background(), "Delete", key, flag); err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		if flag.Flag {
			w.Write([]byte("1"))
		} else {
			w.Write([]byte("0"))
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", router)
	rpc.StartXServer(10086)
	fmt.Println("SDCS is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
