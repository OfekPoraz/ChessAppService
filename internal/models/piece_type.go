package models

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
