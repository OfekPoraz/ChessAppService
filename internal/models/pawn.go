package models

import "fmt"

type Pawn struct {
	Color string
}

func (p Pawn) GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition {
	moves := []PossibleMovesPosition{}
	direction := -1
	if p.Color == "black" {
		direction = 1
	}
	if board.IsValidPosition(PossibleMovesPosition{Row: position.Row + direction, Col: position.Col}) {
		moves = append(moves, PossibleMovesPosition{Row: position.Row + direction, Col: position.Col, IsCapture: false})
	}

	eatMoves := []Moves{{1, 1}, {1, -1}}
	if p.Color == "white" {
		eatMoves = []Moves{{Row: -eatMoves[0].Row, Col: eatMoves[0].Col}, {Row: -eatMoves[1].Row, Col: eatMoves[1].Col}}
	}

	for _, move := range eatMoves {
		newPos := PossibleMovesPosition{Row: position.Row + move.Row, Col: position.Col + move.Col}
		fmt.Println("possible eating pawn", newPos)
		if board.IsValidPosition(newPos) {
			piece := board.GetPiece(newPos)
			if piece != nil && piece.Color != p.Color {
				fmt.Println("eat true")
				moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
			}
		}
	}

	return moves
}

func (p Pawn) GetColor() string {
	return p.Color
}

func (p Pawn) GetType() PieceType {
	return PawnType
}
