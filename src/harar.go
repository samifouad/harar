package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Initialize etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Fetch the HTML file from etcd
		resp, err := cli.Get(ctx, "/path/to/htmlfile")
		if err != nil {
			http.Error(w, "Failed to fetch data from etcd", http.StatusInternalServerError)
			return
		}
		if len(resp.Kvs) == 0 {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Decode the base64 content
		base64Content := resp.Kvs[0].Value
		decodedContent, err := base64.StdEncoding.DecodeString(string(base64Content))
		if err != nil {
			http.Error(w, "Failed to decode base64 content", http.StatusInternalServerError)
			return
		}

		// Set the content type and write the response
		w.Header().Set("Content-Type", "text/html")
		w.Write(decodedContent)
	})

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
