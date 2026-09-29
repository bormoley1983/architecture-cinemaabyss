package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

// config holds runtime configuration for the events service, loaded from
// environment variables at startup.
type config struct {
	Port        string
	KafkaBrokers string
}

// loadConfig reads configuration from environment variables with sensible
// defaults so the service can also run outside of docker-compose.
func loadConfig() config {
	cfg := config{
		Port:         "8082",
		KafkaBrokers: "localhost:9092",
	}
	if v := os.Getenv("PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		cfg.KafkaBrokers = v
	}
	return cfg
}

// Event is the unified envelope written to Kafka. The payload depends on the
// event type (movie / user / payment).
type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// EventResponse is returned to the API caller after an event has been produced.
type EventResponse struct {
	Status string `json:"status"`
	Event  Event  `json:"event"`
}

// Error is the standard error body.
type Error struct {
	Error string `json:"error"`
}

// topicForType maps an event type to its Kafka topic.
func topicForType(eventType string) string {
	switch eventType {
	case "movie":
		return "movie-events"
	case "user":
		return "user-events"
	case "payment":
		return "payment-events"
	default:
		return eventType + "-events"
	}
}

// produceEvent writes a single event to its Kafka topic.
func produceEvent(w *kafka.Writer, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	topic := topicForType(event.Type)
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(event.ID),
		Value: body,
	}

	if err := w.WriteMessages(context.Background(), msg); err != nil {
		return fmt.Errorf("write to kafka: %w", err)
	}

	log.Printf("PRODUCED event type=%s id=%s topic=%s", event.Type, event.ID, topic)
	return nil
}

// startConsumers launches one background consumer per topic. Each consumer
// belongs to the "events-service" group and logs every message it reads, which
// demonstrates that the service both produces and consumes its own events.
func startConsumers(brokers, groupID string) {
	topics := []string{"movie-events", "user-events", "payment-events"}

	for _, topic := range topics {
		go consumeLoop(brokers, groupID, topic)
	}
}

func consumeLoop(brokers, groupID, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokers},
		Topic:     topic,
		GroupID:   groupID,
		MinBytes:  1, // default
		MaxBytes:  10e6,
	})

	log.Printf("CONSUMER started for topic=%s group=%s", topic, groupID)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			// Transient errors (e.g. broker not ready yet) are logged and we retry.
			log.Printf("CONSUMER topic=%s read error: %v (retrying)", topic, err)
			time.Sleep(time.Second)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("CONSUMER topic=%s FAILED to parse message: %v", topic, err)
			continue
		}

		payloadStr := string(msg.Value)
		log.Printf("CONSUMED event from topic=%s partition=%d offset=%d | type=%s id=%s payload=%s",
			topic, msg.Partition, msg.Offset, event.Type, event.ID, payloadStr)
	}
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// newEventHandler builds a POST handler that decodes the payload into `in`,
// produces an event of the given type, and returns the EventResponse.
func newEventHandler(w *kafka.Writer, eventType string, decode func(*http.Request) (interface{}, error)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(rw, http.StatusMethodNotAllowed, Error{Error: "method not allowed"})
			return
		}

		payload, err := decode(r)
		if err != nil {
			writeJSON(rw, http.StatusBadRequest, Error{Error: err.Error()})
			return
		}

		event := Event{
			ID:        uuid.NewString(),
			Type:      eventType,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
		}

		if err := produceEvent(w, event); err != nil {
			log.Printf("ERROR producing %s event: %v", eventType, err)
			writeJSON(rw, http.StatusInternalServerError, Error{Error: err.Error()})
			return
		}

		writeJSON(rw, http.StatusCreated, EventResponse{
			Status: "success",
			Event:  event,
		})
	}
}

func main() {
	cfg := loadConfig()

	log.Printf("Starting events service on port %s", cfg.Port)
	log.Printf("  KAFKA_BROKERS=%s", cfg.KafkaBrokers)

	// Kafka writer (producer).
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaBrokers),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	defer writer.Close()

	// Start background consumers for all three topics.
	startConsumers(cfg.KafkaBrokers, "events-service")

	// Health check.
	http.HandleFunc("/api/events/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"status": true})
	})

	// Event creation endpoints.
	http.HandleFunc("/api/events/movie", newEventHandler(writer, "movie", decodeMovieEvent))
	http.HandleFunc("/api/events/user", newEventHandler(writer, "user", decodeUserEvent))
	http.HandleFunc("/api/events/payment", newEventHandler(writer, "payment", decodePaymentEvent))

	log.Printf("Events service listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

// --- payload decoders (match the API specification schemas) ---

type MovieEvent struct {
	MovieID       int      `json:"movie_id"`
	Title         string   `json:"title"`
	Action        string   `json:"action"`
	UserID        int      `json:"user_id,omitempty"`
	Rating        float64  `json:"rating,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	Description   string   `json:"description,omitempty"`
}

type UserEvent struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentID  int       `json:"payment_id"`
	UserID     int       `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	MethodType string    `json:"method_type,omitempty"`
}

func decodeMovieEvent(r *http.Request) (interface{}, error) {
	var e MovieEvent
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}
	if e.MovieID == 0 || e.Title == "" || e.Action == "" {
		return nil, fmt.Errorf("movie_id, title and action are required")
	}
	return e, nil
}

func decodeUserEvent(r *http.Request) (interface{}, error) {
	var e UserEvent
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}
	if e.UserID == 0 || e.Action == "" {
		return nil, fmt.Errorf("user_id and action are required")
	}
	return e, nil
}

func decodePaymentEvent(r *http.Request) (interface{}, error) {
	var e PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}
	if e.PaymentID == 0 || e.UserID == 0 || e.Status == "" {
		return nil, fmt.Errorf("payment_id, user_id and status are required")
	}
	return e, nil
}
