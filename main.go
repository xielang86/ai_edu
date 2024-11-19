package main

import (
	"api_server/api"
	"fmt"
	"net/http"
)

func main() {
	// http.HandleFunc("/", api.QueryHandler)
	http.HandleFunc("/register", api.RegisterHandler)
	http.HandleFunc("/login", api.LoginHandler)
	fmt.Println("Server starting on port :8080...")
	http.ListenAndServe(":8080", nil)
}
