package pieces

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/models"
	"fmt"
)

type Knight struct {
	core.BasePiece
}

func NewKnight(color string) core.Piece {
	return &Knight{BasePiece: core.BasePiece{Color: color, Type: core.KnightType}}
}

func (k Knight) GetPossibleMoves(board *core.Board, position core.PossibleMovesPosition) []core.PossibleMovesPosition {
	var moves []core.PossibleMovesPosition
	knightMoves := []models.Moves{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	for _, move := range knightMoves {
		newPos := core.PossibleMovesPosition{Row: position.Row + move.Row, Col: position.Col + move.Col}
		fmt.Println(newPos)
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

func (k Knight) GetColor() string {
	return k.Color
}

func (k Knight) GetType() core.PieceType {
	return core.KnightType
}
