package models

type King struct {
	Color string
}

func (k King) GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition {
	moves := []PossibleMovesPosition{}
	kingMoves := []Moves{
		{1, 0}, {1, 1}, {1, -1},
		{0, 1}, {0, -1},
		{-1, 0}, {-1, 1}, {-1, -1},
	}
	for _, move := range kingMoves {
		newPos := PossibleMovesPosition{Row: position.Row + move.Row, Col: position.Col + move.Col}
		if board.IsValidPosition(newPos) {
			if board.GetPiece(newPos) == nil || board.GetPiece(newPos).Color == k.Color {
				moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: false})
			} else {
				moves = append(moves, PossibleMovesPosition{Row: newPos.Row, Col: newPos.Col, IsCapture: true})
			}
		}
	}
	return moves
}

func (k King) GetColor() string {
	return k.Color
}

func (k King) GetType() PieceType {
	return KingType
}
