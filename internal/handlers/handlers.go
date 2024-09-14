package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ChessAppIdoBack/internal/models"
	"ChessAppIdoBack/internal/services"
	"github.com/gorilla/mux"
)

var games = make(map[string]models.Game) // In-memory storage for demonstration

const GameAlreadyExistsErrorCode = "game_already_exists"

func CreateGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the game ID already exists
	if _, exists := games[game.ID]; exists {
		response := map[string]string{
			"error":   GameAlreadyExistsErrorCode,
			"message": "Game with this ID already exists",
		}
		w.WriteHeader(http.StatusConflict) // HTTP 409 Conflict
		json.NewEncoder(w).Encode(response)
		return
	}

	game.ID = "some_unique_id" // Generate unique ID
	game.Board = services.InitializeBoard(7, 5)
	game.CurrentTurn = "white"
	games[game.ID] = game

	response := map[string]interface{}{
		"id":           game.ID,
		"player1_id":   game.Player1ID,
		"player2_id":   game.Player2ID,
		"state":        game.State,
		"current_turn": game.CurrentTurn,
		"board":        game.Board.Pieces,
	}

	fmt.Println(response)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, ok := games[params["id"]]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":           game.ID,
		"player1_id":   game.Player1ID,
		"player2_id":   game.Player2ID,
		"state":        game.State,
		"current_turn": game.CurrentTurn,
		"board":        game.Board.Pieces,
	}

	json.NewEncoder(w).Encode(response)
}

func MakeMove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, ok := games[params["id"]]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	var move struct {
		From models.PossibleMovesPosition `json:"from"`
		To   models.PossibleMovesPosition `json:"to"`
	}
	err := json.NewDecoder(r.Body).Decode(&move)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("move is - ", move.From, move.To)
	if !game.Board.IsValidPosition(move.From) || !game.Board.IsValidPosition(move.To) {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	piece := game.Board.GetPiece(move.From)
	if piece == nil {
		http.Error(w, "No piece at source position", http.StatusBadRequest)
		return
	}

	fmt.Println("piece: ", piece)

	possibleMoves := game.Board.GetPossibleMoves(move.From)
	validMove := false
	fmt.Println("possible moves", possibleMoves)
	for _, pos := range possibleMoves {
		if pos == move.To {
			validMove = true
			break
		}
	}

	fmt.Println("valid move: ", validMove)
	if !validMove {
		http.Error(w, "Invalid move for piece", http.StatusBadRequest)
		return
	}

	isWin, winner := game.Board.MovePiece(game.CurrentTurn, move.From, move.To)
	if isWin {
		response := map[string]interface{}{
			"winner": winner,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	game.CurrentTurn = switchTurn(game.CurrentTurn)
	games[game.ID] = game

	response := map[string]interface{}{
		"id":           game.ID,
		"player1_id":   game.Player1ID,
		"player2_id":   game.Player2ID,
		"state":        game.State,
		"current_turn": game.CurrentTurn,
		"board":        game.Board.Pieces,
	}

	fmt.Println("response - ", response)
	json.NewEncoder(w).Encode(response)
}

func GetGameState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, ok := games[params["id"]]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":           game.ID,
		"player1_id":   game.Player1ID,
		"player2_id":   game.Player2ID,
		"state":        game.State,
		"current_turn": game.CurrentTurn,
		"board":        game.Board.Pieces,
	}

	json.NewEncoder(w).Encode(response)
}

func GetPossibleMoves(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, ok := games[params["id"]]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	piecePosition := r.URL.Query().Get("piecePosition")
	if piecePosition == "" {
		http.Error(w, "piecePosition query parameter is required", http.StatusBadRequest)
		return
	}

	fmt.Println(game.ID)
	fmt.Println(game.CurrentTurn)
	fmt.Println(game.Player2ID)
	fmt.Println(game.Player1ID)
	fmt.Println(game.Board)

	position := models.PositionFromString(piecePosition)
	fmt.Println("PossibleMovesPosition: ")
	fmt.Println(position)
	possibleMoves := game.Board.GetPossibleMoves(position)
	fmt.Println("Possible moves: ")
	fmt.Println(possibleMoves)
	json.NewEncoder(w).Encode(possibleMoves)
}

func switchTurn(currentTurn string) string {
	if currentTurn == "white" {
		return "black"
	}
	return "white"
}

func MakeRandomMove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, ok := games[params["id"]]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	isWin, winner, err := game.Board.MakeRandomMove(game.CurrentTurn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isWin {
		response := map[string]interface{}{
			"winner": winner,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	game.CurrentTurn = switchTurn(game.CurrentTurn)
	games[game.ID] = game

	response := map[string]interface{}{
		"id":           game.ID,
		"player1_id":   game.Player1ID,
		"player2_id":   game.Player2ID,
		"state":        game.State,
		"current_turn": game.CurrentTurn,
		"board":        game.Board.Pieces,
	}

	json.NewEncoder(w).Encode(response)
}

func ApplyPowerBoost0(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]

	game, ok := games[gameID]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	err := game.Board.PowerBoost0(game.CurrentTurn)
	if err != nil {
		http.Error(w, "Failed to apply power boost", http.StatusInternalServerError)
		return
	}

	game.CurrentTurn = switchTurn(game.CurrentTurn)
	games[game.ID] = game

	response := map[string]interface{}{
		"success": true,
	}
	json.NewEncoder(w).Encode(response)
}

func ApplyPowerBoost1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]

	game, ok := games[gameID]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	err := game.Board.PowerBoost1(game.CurrentTurn)
	if err != nil {
		http.Error(w, "Failed to apply power boost", http.StatusInternalServerError)
		return
	}

	game.CurrentTurn = switchTurn(game.CurrentTurn)
	games[game.ID] = game

	response := map[string]interface{}{
		"success": true,
	}
	json.NewEncoder(w).Encode(response)
}

func ApplyPowerBoost2(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	gameID := params["id"]

	game, ok := games[gameID]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	err := game.Board.PowerBoost2(game.CurrentTurn)
	if err != nil {
		http.Error(w, "Failed to apply power boost", http.StatusInternalServerError)
		return
	}

	game.CurrentTurn = switchTurn(game.CurrentTurn)
	games[game.ID] = game

	response := map[string]interface{}{
		"success": true,
	}
	json.NewEncoder(w).Encode(response)
}
