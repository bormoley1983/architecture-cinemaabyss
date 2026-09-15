package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// Database connection
var db *sql.DB

// Models
type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Rating      float64  `json:"rating"`
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Set up HTTP routes
	http.HandleFunc("/api/movies", handleMovies)
	http.HandleFunc("/api/movies/health", handleHealth)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Note: Using a different port than the monolith
	}
	log.Printf("Starting movies microservice on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initDB() {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost/cinemaabyss?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

// Movie handlers
func handleMovies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("id") != "" {
			getMovieByID(w, r)
		} else {
			getAllMovies(w, r)
		}
	case "POST":
		createMovie(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, rating FROM movies")
	fmt.Println("get movies from movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	movies := []Movie{}
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Rating); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get genres for this movie
		genreRows, err := db.Query("SELECT genre FROM movie_genres WHERE movie_id = $1", m.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer genreRows.Close()

		genres := []string{}
		for genreRows.Next() {
			var genre string
			if err := genreRows.Scan(&genre); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			genres = append(genres, genre)
		}
		m.Genres = genres

		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovieByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var m Movie
	err := db.QueryRow("SELECT id, title, description, rating FROM movies WHERE id = $1", id).Scan(&m.ID, &m.Title, &m.Description, &m.Rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get genres for this movie
	genreRows, err := db.Query("SELECT genre FROM movie_genres WHERE movie_id = $1", m.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer genreRows.Close()

	genres := []string{}
	for genreRows.Next() {
		var genre string
		if err := genreRows.Scan(&genre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		genres = append(genres, genre)
	}
	m.Genres = genres

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.QueryRow("INSERT INTO movies (title, description, rating) VALUES ($1, $2, $3) RETURNING id",
		m.Title, m.Description, m.Rating).Scan(&m.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, genre := range m.Genres {
		_, err = tx.Exec("INSERT INTO movie_genres (movie_id, genre) VALUES ($1, $2)", m.ID, genre)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}
