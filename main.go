package main

import (
	"assignment/handlers"
	"fmt"
	"net/http"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

// Create a middleware function to handle panics
func withRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// Setup routes with panic recovery
func setupRoutes() {
	http.HandleFunc("/echo", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMatrix(w, r, "echo")
	}))
	http.HandleFunc("/invert", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMatrix(w, r, "invert")
	}))
	http.HandleFunc("/flatten", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMatrix(w, r, "flatten")
	}))
	http.HandleFunc("/sum", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMatrix(w, r, "sum")
	}))
	http.HandleFunc("/multiply", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMatrix(w, r, "multiply")
	}))
}

func main() {
	setupRoutes()
	fmt.Println("listening to server at 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
