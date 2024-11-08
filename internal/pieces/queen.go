package pieces

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/models"
)

type Queen struct {
	core.BasePiece
}

func NewQueen(color string) core.Piece {
	return &Queen{BasePiece: core.BasePiece{Color: color, Type: core.QueenType}}
}

func (q Queen) GetPossibleMoves(board *core.Board, position core.PossibleMovesPosition) []core.PossibleMovesPosition {
	var moves []core.PossibleMovesPosition
	directions := []models.Moves{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, direction := range directions {
		for i := 1; i < board.Rows; i++ {
			newPos := core.PossibleMovesPosition{Row: position.Row + i*direction.Row, Col: position.Col + i*direction.Col}
			if board.IsValidPosition(newPos) {
				if board.GetPiece(newPos) == nil {
					moves = append(moves, newPos)
				} else {
					if board.GetPiece(newPos).Color == q.Color {
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

func (q Queen) GetColor() string {
	return q.Color
}

func (q Queen) GetType() core.PieceType {
	return core.QueenType
}
