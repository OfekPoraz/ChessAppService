package pieces

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/models"
)

type Bishop struct {
	core.BasePiece
}

func NewBishop(color string) core.Piece {
	return &Bishop{BasePiece: core.BasePiece{Color: color, Type: core.BishopType}}
}

func (b Bishop) GetPossibleMoves(board *core.Board, position core.PossibleMovesPosition) []core.PossibleMovesPosition {
	var moves []core.PossibleMovesPosition
	directions := []models.Moves{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, direction := range directions {
		for i := 1; i < board.Rows; i++ {
			newPos := core.PossibleMovesPosition{Row: position.Row + i*direction.Row, Col: position.Col + i*direction.Col}
			if board.IsValidPosition(newPos) {
				if board.GetPiece(newPos) == nil {
					moves = append(moves, core.PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: false})
				} else if board.GetPiece(newPos).Color == b.Color {
					break
				} else {
					moves = append(moves, core.PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
					break
				}
			} else {
				break
			}
		}
	}
	return moves
}

func (b Bishop) GetColor() string {
	return b.Color
}

func (b Bishop) GetType() core.PieceType {
	return core.BishopType
}
