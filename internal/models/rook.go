package models

type Rook struct {
	Color string
}

func (r Rook) GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition {
	moves := []PossibleMovesPosition{}
	directions := []Moves{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, direction := range directions {
		for i := 1; i < board.Rows; i++ {
			newPos := PossibleMovesPosition{Row: position.Row + i*direction.Row, Col: position.Col + i*direction.Col}
			if board.IsValidPosition(newPos) {
				if board.GetPiece(newPos) == nil {
					moves = append(moves, newPos)
				} else {
					if board.GetPiece(newPos).Color == r.Color {
						break
					} else {
						moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
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

func (r Rook) GetType() PieceType {
	return RookType
}
