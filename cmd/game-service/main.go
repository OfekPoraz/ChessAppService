package main

import (
	"log"
	"net/http"

	"ChessAppIdoBack/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Register routes
	registerGameRoutes(r)

	// Start server
	log.Println("Starting server on :8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// registerGameRoutes sets up all routes related to game handling
func registerGameRoutes(r *mux.Router) {
	// Game Routes
	r.HandleFunc("/games", handlers.CreateGame).Methods("POST")
	r.HandleFunc("/games/{id}", handlers.GetGame).Methods("GET")
	r.HandleFunc("/games/{id}/move", handlers.MakeMove).Methods("POST")
	r.HandleFunc("/games/{id}/state", handlers.GetGameState).Methods("GET")
	r.HandleFunc("/games/{id}/moves", handlers.GetPossibleMoves).Methods("GET")
	r.HandleFunc("/games/{id}/randomMove", handlers.MakeRandomMove).Methods("POST")
	r.HandleFunc("/games/{id}/powerboost/0", handlers.ApplyPowerBoost0).Methods("POST")
	r.HandleFunc("/games/{id}/powerboost/1", handlers.ApplyPowerBoost1).Methods("POST")
	r.HandleFunc("/games/{id}/powerboost/2", handlers.ApplyPowerBoost2).Methods("POST")
	r.HandleFunc("/games/{id}/powerboost/3", handlers.ApplyPowerBoost3).Methods("POST")
}
