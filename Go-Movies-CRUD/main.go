// package main - Every Go program starts with a package declaration
// "main" is special - it tells Go this is an executable program, not a library
package main

// import - This section brings in external code libraries we need
import (
	"encoding/json" // Library to convert Go data to JSON format and vice versa
	"fmt"           // Library for formatted printing (like Printf)
	"log"           // Library for logging errors
	"net/http"      // Library to create web servers and handle HTTP requests

	"github.com/gorilla/mux" // External router library for handling different URL paths
)

// Movie struct - This is like a blueprint/template for movie data
// A struct groups related data together
type Movie struct {
	ID       string    `json:"id"`       // Movie's unique identifier (the `json:"id"` tells Go how to name this in JSON)
	Isbn     string    `json:"isbn"`     // ISBN number (like a book code)
	Title    string    `json:"title"`    // Movie title/name
	Director *Director `json:"director"` // Pointer to Director struct (the * means it's a pointer/reference)
}

// Director struct - Blueprint for director information
type Director struct {
	Firstname string `json:"firstname"` // Director's first name
	Lastname  string `json:"lastname"`  // Director's last name
}

// movies - A global variable that holds a slice (like an array) of Movie objects
// This is our "database" for now (data stored in memory, not a real database)
var movies []Movie

// getMovies - Function that returns ALL movies
// w = http.ResponseWriter (what we send back to the user)
// r = *http.Request (the incoming request from the user)
func getMovies(w http.ResponseWriter, r *http.Request) {
	// Set the response header to tell the browser we're sending JSON data
	w.Header().Set("Content-Type", "application/json")

	// Convert the movies slice to JSON and send it back to the user
	// json.NewEncoder(w) creates a JSON encoder that writes to w
	// .Encode(movies) converts our movies slice to JSON format
	json.NewEncoder(w).Encode(movies)
}

// deleteMovie - Function that deletes a specific movie by ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// mux.Vars(r) extracts URL parameters from the request
	// For example, if URL is /movies/1, params["id"] would be "1"
	params := mux.Vars(r)

	// Loop through all movies with their index and value
	// index = position in the slice (0, 1, 2, etc.)
	// item = the actual movie object at that position
	for index, item := range movies {
		// Check if this movie's ID matches the ID from the URL
		if item.ID == params["id"] {
			// Delete the movie by creating a new slice without this item
			// movies[:index] = everything BEFORE the movie
			// movies[index+1:] = everything AFTER the movie
			// append combines them, effectively removing the movie at 'index'
			movies = append(movies[:index], movies[index+1:]...)
			break // Exit the loop since we found and deleted the movie
		}
	}

	// Send back the updated movies list as JSON
	json.NewEncoder(w).Encode(movies)
}

// getMovie - Function that returns ONE specific movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get the ID from the URL parameters
	params := mux.Vars(r)

	// Loop through movies (we only need the value, not the index, so we use _)
	// _ = blank identifier, means "I don't need this value"
	for _, item := range movies {
		// If we find a movie with matching ID
		if item.ID == params["id"] {
			// Convert this single movie to JSON and send it
			json.NewEncoder(w).Encode(item)
			return // Exit the function immediately (don't continue looping)
		}
	}
	// If no movie found, function ends without sending anything
}

// createMovie - Function that adds a NEW movie to our collection
func createMovie(w http.ResponseWriter, r *http.Request) {
	// Set response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Create an empty Movie variable to store the incoming data
	var movie Movie

	// Read the JSON data from the request body and convert it to a Movie struct
	// r.Body contains the data sent by the user
	// &movie = pointer to our movie variable (so it can be modified)
	// _ = we're ignoring any error (not best practice, but common in examples)
	_ = json.NewDecoder(r.Body).Decode(&movie)

	// Add this new movie to our movies slice
	movies = append(movies, movie)

	// Send back the newly created movie as confirmation
	json.NewEncoder(w).Encode(movie)
}

// updateMovie - Function that updates an existing movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	// Set content type to json (this comment was already in your code)
	w.Header().Set("Content-Type", "application/json")

	// Get the movie ID from the URL
	params := mux.Vars(r)

	// Loop through all movies with index and value
	for index, item := range movies {
		// Find the movie that matches the ID
		if item.ID == params["id"] {
			// FIRST: Delete the old movie (same technique as deleteMovie)
			movies = append(movies[:index], movies[index+1:]...)

			// SECOND: Create a variable for the updated movie data
			var movie Movie

			// Read the new movie data from the request body
			_ = json.NewDecoder(r.Body).Decode(&movie)

			// Make sure the ID stays the same (use the ID from URL, not from the body)
			movie.ID = params["id"]

			// THIRD: Add the updated movie back to the slice
			movies = append(movies, movie)

			// Send back the updated movie as confirmation
			json.NewEncoder(w).Encode(movie)
			return // Exit the function
		}
	}
	// If no movie found with that ID, function ends without doing anything
}

// main - The starting point of our program (every Go program needs a main function)
func main() {
	// Create a new router using gorilla/mux
	// A router decides which function to call based on the URL
	r := mux.NewRouter()

	// Mock Data - Creating some sample movies to work with
	// @todo - implement DB means "in the future, use a real database"
	// &Director{...} creates a new Director and returns a pointer to it
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "454555", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	// Register routes - tell the router which function to call for each URL pattern

	// When someone visits /movies with GET method, call getMovies function
	r.HandleFunc("/movies", getMovies).Methods("GET")

	// When someone visits /movies/1 (or any ID) with GET method, call getMovie
	// {id} is a placeholder that captures any value (like 1, 2, abc, etc.)
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")

	// When someone sends POST request to /movies, call createMovie
	r.HandleFunc("/movies", createMovie).Methods("POST")

	// When someone sends PUT request to /movies/1, call updateMovie
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")

	// When someone sends DELETE request to /movies/1, call deleteMovie
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Print a message to the console so we know the server started
	fmt.Printf("Starting server at port 8000\n")

	// Start the web server on port 8000
	// ":8000" means listen on all network interfaces at port 8000
	// r is our router that handles all the routes
	// if err := ... checks if starting the server caused an error
	if err := http.ListenAndServe(":8000", r); err != nil {
		// If there's an error, log it and stop the program
		log.Fatal(err)
	}
}
