package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if present
	godotenv.Load()

	// Create the router
	r := mux.NewRouter()
	
	// Setup API routes
	apiRouter := r.PathPrefix("/api").Subrouter()
	
	// Register artifact routes
	artifactRouter := apiRouter.PathPrefix("/artifacts").Subrouter()
	artifactRouter.HandleFunc("", listArtifacts).Methods("GET")
	artifactRouter.HandleFunc("", uploadArtifact).Methods("POST")
	artifactRouter.HandleFunc("/{id}", getArtifact).Methods("GET")
	artifactRouter.HandleFunc("/{id}", deleteArtifact).Methods("DELETE")
	artifactRouter.HandleFunc("/{id}/download", downloadArtifact).Methods("GET")
	
	// Register authentication routes
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", login).Methods("POST")
	authRouter.HandleFunc("/logout", logout).Methods("POST")
	
	// Middleware for authentication and logging
	apiRouter.Use(loggingMiddleware)
	
	// Define the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// Authentication handler functions
func login(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user authentication
	w.WriteHeader(http.StatusNotImplemented)
}

func logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user logout
	w.WriteHeader(http.StatusNotImplemented)
}

// Artifact handler functions
func listArtifacts(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement artifact listing
	w.WriteHeader(http.StatusNotImplemented)
}

func uploadArtifact(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement artifact upload
	w.WriteHeader(http.StatusNotImplemented)
}

func getArtifact(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement artifact retrieval
	w.WriteHeader(http.StatusNotImplemented)
}

func deleteArtifact(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement artifact deletion
	w.WriteHeader(http.StatusNotImplemented)
}

func downloadArtifact(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement artifact download
	w.WriteHeader(http.StatusNotImplemented)
}