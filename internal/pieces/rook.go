package pieces

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/models"
)

type Rook struct {
	core.BasePiece
}

func NewRook(color string) core.Piece {
	return &Rook{BasePiece: core.BasePiece{Color: color, Type: core.RookType}}
}

func (r Rook) GetPossibleMoves(board *core.Board, position core.PossibleMovesPosition) []core.PossibleMovesPosition {
	var moves []core.PossibleMovesPosition
	directions := []models.Moves{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, direction := range directions {
		for i := 1; i < board.Rows; i++ {
			newPos := core.PossibleMovesPosition{Row: position.Row + i*direction.Row, Col: position.Col + i*direction.Col}
			if board.IsValidPosition(newPos) {
				if board.GetPiece(newPos) == nil {
					moves = append(moves, newPos)
				} else {
					if board.GetPiece(newPos).Color == r.Color {
						break
					} else {
						moves = append(moves, core.PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
						break
					}
				}
			}
		}
	}
	return moves
}

func (r Rook) GetColor() string {
	return r.Color
}

func (r Rook) GetType() core.PieceType {
	return core.RookType
}
