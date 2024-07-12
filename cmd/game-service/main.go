package main

import (
	"log"
	"net/http"

	"ChessAppIdoBack/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Game Routes
	r.HandleFunc("/games", handlers.CreateGame).Methods("POST")
	r.HandleFunc("/games/{id}", handlers.GetGame).Methods("GET")
	r.HandleFunc("/games/{id}/move", handlers.MakeMove).Methods("POST")
	r.HandleFunc("/games/{id}/state", handlers.GetGameState).Methods("GET")
	r.HandleFunc("/games/{id}/moves", handlers.GetPossibleMoves).Methods("GET") // New endpoint

	log.Fatal(http.ListenAndServe(":8000", r))
}
