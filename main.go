package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"rpc"
)

var neighbors []string
var xport = 10086
var httpport = 8080

func router(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	xclient := rpc.NewXClient("localhost", xport)
	defer xclient.Close()

	switch r.Method {
	case http.MethodGet:
		// Prepare args
		key := rpc.Key{Key: r.URL.Path[1:]}
		value := &rpc.Value{}

		// Search local
		if err := xclient.Call(context.Background(), "Get", key, value); err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		if value.Value != "" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value.Value))
			return
		}

		// Search neighbors
		for _, neighbor := range neighbors {
			xclient_neighbor := rpc.NewXClient(neighbor, xport)
			defer xclient_neighbor.Close()
			if err := xclient_neighbor.Call(context.Background(), "Get", key, value); err != nil {
				log.Fatalf("failed to call: %v", err)
			}
			if value.Value != "" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(value.Value))
				return
			}
		}

		// Not found
		http.Error(w, "Not found", http.StatusNotFound)
	case http.MethodPost:
		content := string(body)

		// Extract key
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

		// Prepare args
		key := content[key_start+1 : key_start+1+key_end]
		pair := rpc.Pair{
			Key:   key,
			Value: content,
		}
		flag := &rpc.Flag{}

		// Update local
		if err := xclient.Call(context.Background(), "Post", pair, flag); err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		// Update neighbors
		for _, neighbor := range neighbors {
			xclient_neighbor := rpc.NewXClient(neighbor, xport)
			defer xclient_neighbor.Close()

			// Query whether need to update the value
			if err := xclient.Call(context.Background(), "Query", pair, flag); err != nil {
				log.Fatalf("failed to call: %v", err)
			}
			if !flag.Flag {
				continue
			}
			if err := xclient_neighbor.Call(context.Background(), "Post", pair, flag); err != nil {
				log.Fatalf("failed to call: %v", err)
			}
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		result := false

		// Prepare args
		key := rpc.Key{Key: r.URL.Path[1:]}
		flag := &rpc.Flag{}

		// Delete local
		if err := xclient.Call(context.Background(), "Delete", key, flag); err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		if flag.Flag {
			result = true
		}

		// Delete neighbors
		for _, neighbor := range neighbors {
			xclient_neighbor := rpc.NewXClient(neighbor, xport)
			defer xclient_neighbor.Close()
			if err := xclient_neighbor.Call(context.Background(), "Delete", key, flag); err != nil {
				log.Fatalf("failed to call: %v", err)
			}
			if flag.Flag {
				result = true
			}
		}

		// Reply
		if result {
			w.Write([]byte("1"))
		} else {
			w.Write([]byte("0"))
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	for i := 1; i < len(os.Args); i++ {
		neighbors = append(neighbors, os.Args[i])
	}
	fmt.Printf("Neighbors: %v\n", neighbors)

	rpc.StartXServer(xport)

	http.HandleFunc("/", router)
	fmt.Printf("SDCS is running on http://localhost:%d\n", httpport)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpport), nil))
}
