package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {

	http.HandleFunc("/", DefaultHandler)
	http.HandleFunc("/variety", LoggerMiddleware(http.HandlerFunc(VarietyHandler)))

	fmt.Println("> starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}

// LoggerMiddleware ...
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)

	})

}

// DefaultHandler ...
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello to /")
	
}

// VarietyHandler ...
func VarietyHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello to /variety")

}