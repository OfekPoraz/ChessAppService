package services

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/pieces"
)

func InitializeBoard(rows, columns int) *core.Board {
	board := core.NewBoard(rows, columns, pieces.CreatePiece)

	// Initialize white and black pieces using helper function
	initializeSide(board, "white", 1, 0)
	initializeSide(board, "black", 5, 6)

	board.PlaceMine("white")
	return board
}

// initializeSide initializes all pieces for one side (white or black)
func initializeSide(board *core.Board, color string, pawnRow, mainRow int) {
	// Place pawns
	for col := 0; col < 5; col++ {
		board.PlacePiece(pieces.NewPawn(color), core.PossibleMovesPosition{Row: pawnRow, Col: col})
	}

	// Place main pieces
	piecePlacement := []struct {
		Piece core.Piece
		Col   int
	}{
		{pieces.NewRook(color), 0},
		{pieces.NewKing(color), 1},
		{pieces.NewKnight(color), 2},
		{pieces.NewQueen(color), 3},
		{pieces.NewBishop(color), 4},
	}

	for _, pieceData := range piecePlacement {
		board.PlacePiece(pieceData.Piece, core.PossibleMovesPosition{Row: mainRow, Col: pieceData.Col})
	}
}
