package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sudarakas/edata/service/subscription"
	"github.com/sudarakas/edata/service/user"
)

// Config holds server configuration
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// APIServer represents the API server
type APIServer struct {
	config Config
	db     *sql.DB
	router *mux.Router
	server *http.Server
}

// NewAPIServer creates a new instance of APIServer
func NewAPIServer(config Config, db *sql.DB) *APIServer {
	if db == nil {
		panic("database connection cannot be nil")
	}

	router := mux.NewRouter()

	return &APIServer{
		config: config,
		db:     db,
		router: router,
		server: &http.Server{
			Addr:         config.Addr,
			Handler:      securityHeadersMiddleware(router),
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		},
	}
}

// securityHeadersMiddleware adds common security headers to each response.
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none';")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

// setupRoutes initializes all route handlers
func (s *APIServer) setupRoutes() error {
	subRouter := s.router.PathPrefix("/api/v1").Subrouter()

	userStore, err := user.NewStore(s.db)
	if err != nil {
		return fmt.Errorf("failed to create user store: %w", err)
	}

	subscriptionStore, err := subscription.NewStore(s.db)
	if err != nil {
		return fmt.Errorf("failed to create subscription store: %w", err)
	}

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoute(subRouter)

	subscriptionHandler := subscription.NewHandler(subscriptionStore)
	subscriptionHandler.RegisterRoutes(subRouter)

	s.router.HandleFunc("/health", s.healthCheck()).Methods("GET")

	return nil
}

// Run starts the server and handles graceful shutdown
func (s *APIServer) Run() error {
	if err := s.setupRoutes(); err != nil {
		return fmt.Errorf("failed to setup routes: %w", err)
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("API listening on %s", s.config.Addr)
		serverErrors <- s.server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("start shutdown due to %v signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := s.server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			s.server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// HealthCheck returns a simple health check handler
func (s *APIServer) healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// DefaultConfig returns a default server configuration
func DefaultConfig() Config {
	return Config{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
