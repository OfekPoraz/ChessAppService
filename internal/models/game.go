package models

import "ChessAppIdoBack/internal/core"

type Game struct {
	ID          string      `json:"id"`
	Player1ID   string      `json:"player1_id"`
	Player2ID   string      `json:"player2_id"`
	State       string      `json:"state"`
	CurrentTurn string      `json:"current_turn"`
	Board       *core.Board `json:"board"`
}
