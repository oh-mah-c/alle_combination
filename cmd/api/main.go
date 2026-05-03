package main

import (
	"alle_combination/internal/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handler.CalculatePunnet)
	http.HandleFunc("/api/v1/linked-calculate", handler.CalculateLinkedPunnet)

	port := ":8080"

	fmt.Println("Mendel Genetics is activating")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
