package main

import (
	"dopc/internal/controller"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/delivery-order-price", controller.DopcController)
	fmt.Println("Testserver up at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
