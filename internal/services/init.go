package services

import "ChessAppIdoBack/internal/models"

func InitializeBoard(rows, columns int) *models.Board {
	board := models.NewBoard(rows, columns)

	// Place white pieces (only the left 5 columns)
	for col := 0; col < 5; col++ {
		board.PlacePiece(models.Pawn{Color: "white"}, models.PossibleMovesPosition{Row: 1, Col: col})
	}
	board.PlacePiece(models.Rook{Color: "white"}, models.PossibleMovesPosition{Row: 0, Col: 0})
	board.PlacePiece(models.Knight{Color: "white"}, models.PossibleMovesPosition{Row: 0, Col: 1})
	board.PlacePiece(models.Bishop{Color: "white"}, models.PossibleMovesPosition{Row: 0, Col: 2})
	board.PlacePiece(models.Queen{Color: "white"}, models.PossibleMovesPosition{Row: 0, Col: 3})
	board.PlacePiece(models.King{Color: "white"}, models.PossibleMovesPosition{Row: 0, Col: 4})

	// Place black pieces (only the left 5 columns)
	for col := 0; col < 5; col++ {
		board.PlacePiece(models.Pawn{Color: "black"}, models.PossibleMovesPosition{Row: 5, Col: col})
	}
	board.PlacePiece(models.Rook{Color: "black"}, models.PossibleMovesPosition{Row: 6, Col: 0})
	board.PlacePiece(models.Knight{Color: "black"}, models.PossibleMovesPosition{Row: 6, Col: 1})
	board.PlacePiece(models.Bishop{Color: "black"}, models.PossibleMovesPosition{Row: 6, Col: 2})
	board.PlacePiece(models.Queen{Color: "black"}, models.PossibleMovesPosition{Row: 6, Col: 3})
	board.PlacePiece(models.King{Color: "black"}, models.PossibleMovesPosition{Row: 6, Col: 4})

	return board
}
