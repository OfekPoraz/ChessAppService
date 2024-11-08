package pieces

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/models"
)

type King struct {
	core.BasePiece
}

func NewKing(color string) core.Piece {
	return &King{BasePiece: core.BasePiece{Color: color, Type: core.KingType}}
}

func (k King) GetPossibleMoves(board *core.Board, position core.PossibleMovesPosition) []core.PossibleMovesPosition {
	var moves []core.PossibleMovesPosition
	kingMoves := []models.Moves{
		{1, 0}, {1, 1}, {1, -1},
		{0, 1}, {0, -1},
		{-1, 0}, {-1, 1}, {-1, -1},
	}
	for _, move := range kingMoves {
		newPos := core.PossibleMovesPosition{Row: position.Row + move.Row, Col: position.Col + move.Col}
		if board.IsValidPosition(newPos) {
			if board.GetPiece(newPos) == nil {
				moves = append(moves, core.PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: false})
			} else if board.GetPiece(newPos).Color != k.Color {
				moves = append(moves, core.PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
			}
		}
	}
	return moves
}

func (k King) GetColor() string {
	return k.Color
}

func (k King) GetType() core.PieceType {
	return core.KingType
}
