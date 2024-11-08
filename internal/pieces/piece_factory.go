package pieces

import (
	"ChessAppIdoBack/internal/core"
)

var pieceConstructors = map[core.PieceType]func(string) core.Piece{
	core.PawnType:   func(color string) core.Piece { return NewPawn(color) },
	core.RookType:   func(color string) core.Piece { return NewRook(color) },
	core.KnightType: func(color string) core.Piece { return NewKnight(color) },
	core.BishopType: func(color string) core.Piece { return NewBishop(color) },
	core.QueenType:  func(color string) core.Piece { return NewQueen(color) },
	core.KingType:   func(color string) core.Piece { return NewKing(color) },
}

func CreatePiece(color string, pieceType core.PieceType) core.Piece {
	if constructor, found := pieceConstructors[pieceType]; found {
		return constructor(color)
	}
	return nil
}
