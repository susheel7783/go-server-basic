# Go Movies CRUD API

A simple REST API built with Go and Gorilla Mux for managing movies.

## Features
- Get all movies
- Get a single movie by ID
- Create a new movie
- Update an existing movie
- Delete a movie

## Installation
```bash
go get github.com/gorilla/mux
```

## Run
```bash
go run main.go
```

Server will start at `http://localhost:8000`

## API Endpoints

- `GET /movies` - Get all movies
- `GET /movies/{id}` - Get a movie by ID
- `POST /movies` - Create a new movie
- `PUT /movies/{id}` - Update a movie
- `DELETE /movies/{id}` - Delete a movie
