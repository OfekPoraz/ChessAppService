package models

type Bishop struct {
	Color string
}

func (b Bishop) GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition {
	moves := []PossibleMovesPosition{}
	directions := []Moves{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, direction := range directions {
		for i := 1; i < board.Rows; i++ {
			newPos := PossibleMovesPosition{Row: position.Row + i*direction.Row, Col: position.Col + i*direction.Col}
			if board.IsValidPosition(newPos) {
				if board.GetPiece(newPos) == nil {
					moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: false})
				} else if board.GetPiece(newPos).Color == b.Color {
					break
				} else {
					moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
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

func (b Bishop) GetType() PieceType {
	return BishopType
}
