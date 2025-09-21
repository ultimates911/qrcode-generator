package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		type resp struct{ Message string `json:"message"` }
		_ = json.NewEncoder(w).Encode(resp{Message: "привет из Go"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
