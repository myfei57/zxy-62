package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tunnelnet/internal/console"
	"tunnelnet/internal/store"
)

func main() {
	webDir := os.Getenv("TUNNELNET_WEB_DIR")
	if webDir == "" {
		webDir = "web"
	}
	dataDir := os.Getenv("TUNNELNET_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	storage, err := store.New(dataDir)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	server := console.New(webDir, storage)
	addr := os.Getenv("TUNNELNET_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	fmt.Printf("listening on %s\n", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
