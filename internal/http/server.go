package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-bot/internal/logging"
)

// Server represents an HTTP server for health checks
type Server struct {
	server *http.Server
}

// HealthResponse is the JSON structure for health check responses
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// NewServer creates a new HTTP server for health checks
func NewServer(port int) *Server {
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/health", handleHealth)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &Server{
		server: server,
	}
}

// Start starts the HTTP server
func (s *Server) Start() {
	go func() {
		logging.Info("Starting health check server on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Health check server error: %v", err)
		}
	}()
}

// Stop gracefully shuts down the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	logging.Info("Shutting down health check server")
	return s.server.Shutdown(ctx)
}

// Health check handler
func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "0.1.0", // Update with your actual version
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
