package models

import "fmt"

type PiecePosition struct {
	Location string    `json:"location"`
	Color    string    `json:"color"`
	Type     PieceType `json:"type"`
}

type Board struct {
	Rows    int             `json:"rows"`
	Columns int             `json:"columns"`
	Pieces  []PiecePosition `json:"pieces"`
}

func NewBoard(rows, columns int) *Board {
	return &Board{
		Rows:    rows,
		Columns: columns,
		Pieces:  []PiecePosition{},
	}
}

func (b *Board) PlacePiece(piece Piece, position PossibleMovesPosition) {
	if position.Row >= 0 && position.Row < b.Rows && position.Col >= 0 && position.Col < b.Columns {
		location := fmt.Sprintf("%c%d", 'A'+position.Col, b.Rows-position.Row)
		piecePosition := PiecePosition{
			Location: location,
			Color:    piece.GetColor(),
			Type:     piece.GetType(),
		}
		b.Pieces = append(b.Pieces, piecePosition)
	} else {
		fmt.Printf("Invalid position for placing piece: Row=%d, Col=%d\n", position.Row, position.Col)
		panic("PossibleMovesPosition out of bounds")
	}
}

func (b *Board) MovePiece(from, to PossibleMovesPosition) {
	fromLocation := fmt.Sprintf("%c%d", 'A'+from.Col, from.Row+1)
	fmt.Println("Moving from", fromLocation)
	toLocation := fmt.Sprintf("%c%d", 'A'+to.Col, to.Row+1)
	fmt.Println("Moving to", toLocation)
	// Remove any piece that is currently at the destination location
	for i, piece := range b.Pieces {
		if piece.Location == toLocation {
			// Remove the piece at the destination
			b.Pieces = append(b.Pieces[:i], b.Pieces[i+1:]...)
			break
		}
	}

	// Move the piece from the source location to the destination
	for i, piece := range b.Pieces {
		if piece.Location == fromLocation {
			b.Pieces[i].Location = toLocation
			return
		}
	}
}

func (b *Board) GetPiece(position PossibleMovesPosition) *PiecePosition {
	location := fmt.Sprintf("%c%d", 'A'+position.Col, position.Row+1)
	for _, piece := range b.Pieces {
		if piece.Location == location {
			return &piece
		}
	}
	return nil
}

func (b *Board) IsValidPosition(position PossibleMovesPosition) bool {
	return position.Row >= 0 && position.Row < b.Rows && position.Col >= 0 && position.Col < b.Columns
}

func (b *Board) GetPossibleMoves(position PossibleMovesPosition) []PossibleMovesPosition {
	piecePos := b.GetPiece(position)
	fmt.Println("Piece is:")
	fmt.Println(piecePos)
	if piecePos == nil {
		return nil
	}
	piece := CreatePiece(piecePos.Color, piecePos.Type)
	fmt.Println("Piece:")
	fmt.Println(piece.GetColor(), piece.GetType())
	return piece.GetPossibleMoves(b, position)
}

func PositionFromString(posStr string) PossibleMovesPosition {
	col := int(posStr[0] - 'A')
	row := int(posStr[1] - '1')
	return PossibleMovesPosition{Row: row, Col: col}
}
