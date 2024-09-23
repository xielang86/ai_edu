package main

import (
	"api_server/api"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", api.QueryHandler)
	fmt.Println("Server starting on port :8080...")
	http.ListenAndServe(":8080", nil)
}
