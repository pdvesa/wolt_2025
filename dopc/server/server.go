package server

import (
	"dopc/internal/controller"
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/api/v1/delivery-order-price", controller.ApiController)
	fmt.Println("Testserver up at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
