// ==============================
// Package Declaration
// ==============================

// package main tells Go that this file is an executable program
// Every runnable Go program MUST have package main
package main

// ==============================
// Importing Libraries
// ==============================

import (
	// fmt is used to format and print text
	// In web servers, fmt.Fprintf is used to send text to the browser
	"fmt"

	// log is used to print error messages in the terminal
	// log.Fatal() prints the error and stops the program
	"log"

	// net/http is the core library for building web servers in Go
	// It provides Request, ResponseWriter, routing, and server tools
	"net/http"
)

// ==============================
// formHandler Function
// ==============================

// formHandler is a HANDLER FUNCTION
// A handler runs when a specific URL is called (here: /form)
//
// w (ResponseWriter) → used to SEND response back to browser
// r (*Request)       → contains data SENT by the browser
func formHandler(w http.ResponseWriter, r *http.Request) {

	// ParseForm reads form data from the HTTP request body
	// Without this, FormValue() may not work correctly
	if err := r.ParseForm(); err != nil {

		// If parsing fails, send error message to browser
		// fmt.Fprintf writes data INTO w (response)
		fmt.Fprintf(w, "ParseForm() err: %v", err)

		// return stops execution of this function
		return
	}

	// If parsing is successful, send confirmation to browser
	fmt.Fprintf(w, "POST request successful\n")

	// Read the value of input field named "name"
	// Example: <input name="name">
	name := r.FormValue("name")

	// Read the value of input field named "address"
	address := r.FormValue("address")

	// Send the name value back to browser
	// %s is a format specifier for string
	fmt.Fprintf(w, "Name = %s\n", name)

	// Send the address value back to browser
	fmt.Fprintf(w, "Address = %s\n", address)
}

// ==============================
// helloHandler Function
// ==============================

// helloHandler handles requests to /hello
// This handler only allows GET requests
func helloHandler(w http.ResponseWriter, r *http.Request) {

	// r.URL.Path contains the URL path requested by browser
	// This ensures only "/hello" is accepted
	if r.URL.Path != "/hello" {

		// http.Error sends an HTTP error response to browser
		// StatusNotFound = 404
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	// r.Method contains HTTP method (GET, POST, etc.)
	// This ensures only GET requests are allowed
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}

	// Send response text to browser
	// fmt.Fprintf writes into ResponseWriter (w)
	fmt.Fprintf(w, "hello!")
}

// ==============================
// main Function (Program Start)
// ==============================

// main is the ENTRY POINT of the program
// Go starts execution from here
func main() {

	// Create a file server to serve static files
	// http.Dir("./go_server") tells Go where files are located
	fileServer := http.FileServer(http.Dir("./go_server"))

	// Map "/" (root URL) to the file server
	// Example: http://localhost:8080/
	http.Handle("/", fileServer)

	// Map "/form" URL to formHandler function
	// When /form is called → formHandler runs
	http.HandleFunc("/form", formHandler)

	// Map "/hello" URL to helloHandler function
	http.HandleFunc("/hello", helloHandler)

	// Print message in terminal when server starts
	fmt.Println("Starting server at port 8080")

	// Start HTTP server on port 8080
	// nil means: use Go's default router (ServeMux)
	if err := http.ListenAndServe(":8080", nil); err != nil {

		// If server fails to start, log the error and exit
		log.Fatal(err)
	}
}

// ==============================
// How to run this program
// ==============================
//
// go run go_server/main.go
//
