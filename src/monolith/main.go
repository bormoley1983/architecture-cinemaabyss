package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Database connection
var db *sql.DB

// Models
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Rating      float64  `json:"rating"`
}

type Payment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type Subscription struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PlanType  string    `json:"plan_type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Set up HTTP routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/movies", handleMovies)
	http.HandleFunc("/api/payments", handlePayments)
	http.HandleFunc("/api/subscriptions", handleSubscriptions)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

// User handlers
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("id") != "" {
			getUserByID(w, r)
		} else {
			getAllUsers(w, r)
		}
	case "POST":
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var u User
	err := db.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id", u.Username, u.Email).Scan(&u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
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
	fmt.Println("get movies from monolith")
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

// Payment handlers
func handlePayments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("id") != "" {
			getPaymentByID(w, r)
		} else if r.URL.Query().Get("user_id") != "" {
			getPaymentsByUserID(w, r)
		} else {
			getAllPayments(w, r)
		}
	case "POST":
		createPayment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllPayments(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, user_id, amount, timestamp FROM payments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	payments := []Payment{}
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.Amount, &p.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments = append(payments, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func getPaymentByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var p Payment
	err := db.QueryRow("SELECT id, user_id, amount, timestamp FROM payments WHERE id = $1", id).Scan(&p.ID, &p.UserID, &p.Amount, &p.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func getPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	rows, err := db.Query("SELECT id, user_id, amount, timestamp FROM payments WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	payments := []Payment{}
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.Amount, &p.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments = append(payments, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	var p Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Timestamp = time.Now()
	err := db.QueryRow("INSERT INTO payments (user_id, amount, timestamp) VALUES ($1, $2, $3) RETURNING id",
		p.UserID, p.Amount, p.Timestamp).Scan(&p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// Subscription handlers
func handleSubscriptions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Query().Get("id") != "" {
			getSubscriptionByID(w, r)
		} else if r.URL.Query().Get("user_id") != "" {
			getSubscriptionsByUserID(w, r)
		} else {
			getAllSubscriptions(w, r)
		}
	case "POST":
		createSubscription(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, user_id, plan_type, start_date, end_date FROM subscriptions")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	subscriptions := []Subscription{}
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.PlanType, &s.StartDate, &s.EndDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		subscriptions = append(subscriptions, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

func getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var s Subscription
	err := db.QueryRow("SELECT id, user_id, plan_type, start_date, end_date FROM subscriptions WHERE id = $1", id).Scan(
		&s.ID, &s.UserID, &s.PlanType, &s.StartDate, &s.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func getSubscriptionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	rows, err := db.Query("SELECT id, user_id, plan_type, start_date, end_date FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	subscriptions := []Subscription{}
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.PlanType, &s.StartDate, &s.EndDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		subscriptions = append(subscriptions, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

func createSubscription(w http.ResponseWriter, r *http.Request) {
	var s Subscription
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO subscriptions (user_id, plan_type, start_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id",
		s.UserID, s.PlanType, s.StartDate, s.EndDate).Scan(&s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}
