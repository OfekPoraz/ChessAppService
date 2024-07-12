package models

func CreatePiece(color string, pieceType PieceType) Piece {
	switch pieceType {
	case PawnType:
		return Pawn{Color: color}
	case RookType:
		return Rook{Color: color}
	case KnightType:
		return Knight{Color: color}
	case BishopType:
		return Bishop{Color: color}
	case QueenType:
		return Queen{Color: color}
	case KingType:
		return King{Color: color}
	default:
		return nil
	}
}
