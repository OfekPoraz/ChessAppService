package core

type Piece interface {
	GetPossibleMoves(board *Board, position PossibleMovesPosition) []PossibleMovesPosition
	GetColor() string
	GetType() PieceType
}

type PieceType int

const (
	PawnType PieceType = iota
	RookType
	KnightType
	BishopType
	QueenType
	KingType
)

func (pt PieceType) String() string {
	return [...]string{"pawn", "rook", "knight", "bishop", "queen", "king"}[pt]
}

type BasePiece struct {
	Color string
	Type  PieceType
}

func (bp BasePiece) GetColor() string {
	return bp.Color
}

func (bp BasePiece) GetType() PieceType {
	return bp.Type
}
