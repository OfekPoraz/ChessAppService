package models

import "fmt"

type Knight struct {
	Color string
}

func (k Knight) GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition {
	moves := []PossibleMovesPosition{}
	knightMoves := []Moves{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	for _, move := range knightMoves {
		newPos := PossibleMovesPosition{Row: position.Row + move.Row, Col: position.Col + move.Col}
		fmt.Println(newPos)
		if board.IsValidPosition(newPos) {
			if board.GetPiece(newPos) == nil {
				moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: false})
			} else if board.GetPiece(newPos).Color != k.Color {
				moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
			}
		}
	}
	return moves
}

func (k Knight) GetColor() string {
	return k.Color
}

func (k Knight) GetType() PieceType {
	return KnightType
}
