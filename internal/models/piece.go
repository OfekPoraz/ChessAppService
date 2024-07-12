package models

type Piece interface {
	GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition
	GetColor() string
	GetType() PieceType
}
