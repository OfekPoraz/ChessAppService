package pieces

import (
	"ChessAppIdoBack/internal/core"
	"ChessAppIdoBack/internal/models"
	"fmt"
)

type Pawn struct {
	core.BasePiece
}

func NewPawn(color string) core.Piece {
	return &Pawn{BasePiece: core.BasePiece{Color: color, Type: core.PawnType}}
}

func (p Pawn) GetPossibleMoves(board *core.Board, position core.PossibleMovesPosition) []core.PossibleMovesPosition {
	var moves []core.PossibleMovesPosition
	direction := -1
	if p.Color == "black" {
		direction = 1
	}
	if (board.IsValidPosition(core.PossibleMovesPosition{Row: position.Row + direction, Col: position.Col}) &&
		board.GetPiece(core.PossibleMovesPosition{Row: position.Row + direction, Col: position.Col}) == nil) {
		moves = append(moves, core.PossibleMovesPosition{Row: position.Row + direction, Col: position.Col, IsCapture: false})
	}

	eatMoves := []models.Moves{{1, 1}, {1, -1}}
	if p.Color == "white" {
		eatMoves = []models.Moves{{Row: -eatMoves[0].Row, Col: eatMoves[0].Col}, {Row: -eatMoves[1].Row, Col: eatMoves[1].Col}}
	}

	for _, move := range eatMoves {
		newPos := core.PossibleMovesPosition{Row: position.Row + move.Row, Col: position.Col + move.Col}
		fmt.Println("possible eating pawn", newPos)
		if board.IsValidPosition(newPos) {
			piece := board.GetPiece(newPos)
			if piece != nil && piece.Color != p.Color {
				fmt.Println("eat true")
				moves = append(moves, core.PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
			}
		}
	}

	return moves
}

func (p Pawn) GetColor() string {
	return p.Color
}

func (p Pawn) GetType() core.PieceType {
	return core.PawnType
}
