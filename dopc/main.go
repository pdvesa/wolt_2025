package main

import (
	"dopc/internal/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/delivery-order-price", handler.DopcHandler)
	fmt.Println("Testserver up at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
