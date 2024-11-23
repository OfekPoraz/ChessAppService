package core

import (
	"fmt"
	"math"
	"math/rand"
)

type PiecePosition struct {
	Location string    `json:"location"`
	Color    string    `json:"color"`
	Type     PieceType `json:"type"`
}

type Mine struct {
	Position PossibleMovesPosition `json:"position"`
	Owner    string                `json:"owner"`    // The player who placed the mine
	IsActive bool                  `json:"isActive"` // Mine is active until triggered
}

type PieceFactoryFunc func(color string, pieceType PieceType) Piece

type Board struct {
	Rows         int              `json:"rows"`
	Columns      int              `json:"columns"`
	Pieces       []PiecePosition  `json:"pieces"`
	Mines        []Mine           `json:"mines"` // Added this to track mines on the board
	PieceFactory PieceFactoryFunc // Function to create piece instances
}

var validMoves []struct {
	from PossibleMovesPosition
	to   PossibleMovesPosition
}

func NewBoard(rows, columns int, pieceFactory PieceFactoryFunc) *Board {
	return &Board{
		Rows:         rows,
		Columns:      columns,
		Pieces:       []PiecePosition{},
		Mines:        []Mine{},
		PieceFactory: pieceFactory,
	}
}

func (b *Board) PlacePiece(piece Piece, position PossibleMovesPosition) {
	if position.Row >= 0 && position.Row < b.Rows && position.Col >= 0 && position.Col < b.Columns {
		location := fmt.Sprintf("%c%d", 'A'+position.Col, b.Rows-position.Row)
		fmt.Println("positioned piece", piece, "in position", location)
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

func (b *Board) MovePiece(currentTurn string, from, to PossibleMovesPosition) (bool, string, string) {
	fromLocation := fmt.Sprintf("%c%d", 'A'+from.Col, from.Row+1)
	fmt.Println("Moving from", fromLocation)
	toLocation := fmt.Sprintf("%c%d", 'A'+to.Col, to.Row+1)
	fmt.Println("Moving to", toLocation)

	var capturedPiece *PiecePosition
	var killReason = "" // Field to track the reason for piece elimination

	// Remove any piece that is currently at the destination location
	for i, piece := range b.Pieces {
		if piece.Location == toLocation {
			// Remove the piece at the destination
			b.Pieces = append(b.Pieces[:i], b.Pieces[i+1:]...)
			capturedPiece = &piece
			break
		}
	}

	// Move the piece from the source location to the destination
	for i, piece := range b.Pieces {
		if piece.Location == fromLocation {
			b.Pieces[i].Location = toLocation

			if piece.Type == PawnType && (to.Row == 0 || to.Row == b.Rows-1) {
				// Promote the pawn to a random piece (excluding Pawn and King)
				newPieceType := getRandomPromotionPiece()
				b.Pieces[i].Type = newPieceType
				fmt.Printf("Pawn at %s promoted to %s\n", toLocation, newPieceType)
			}

			break
		}
	}

	if capturedPiece != nil && capturedPiece.Type == KingType {
		var winnerColor string
		if capturedPiece.Color == "white" {
			winnerColor = "black"
		} else {
			winnerColor = "white"
		}
		return true, winnerColor, killReason
	}

	// Check for mines
	mineTriggered, mineMessage := b.TriggerMineIfPresent(to, currentTurn)
	if mineTriggered {
		fmt.Println(mineMessage)
		killReason = "mine"
		// Return the killReason and no win yet
	}

	enemyPieces := []PiecePosition{}
	for _, piece := range b.Pieces {
		if piece.Color != currentTurn {
			enemyPieces = append(enemyPieces, piece)
		}
	}

	// Check if the only remaining enemy piece is the king
	if len(enemyPieces) == 1 && enemyPieces[0].Type == KingType {
		return true, currentTurn, killReason
	}

	currentTurnPieces := []PiecePosition{}
	for _, piece := range b.Pieces {
		if piece.Color == currentTurn {
			currentTurnPieces = append(currentTurnPieces, piece)
		}
	}

	// Check if the only remaining current turn piece is the king
	if len(currentTurnPieces) == 1 && currentTurnPieces[0].Type == KingType {
		var winnerColor string
		if currentTurn == "white" {
			winnerColor = "black"
		} else {
			winnerColor = "white"
		}
		return true, winnerColor, killReason
	}

	return false, "", killReason
}

func (b *Board) MovePieceForChecking(from, to PossibleMovesPosition) (bool, string) {
	fromLocation := fmt.Sprintf("%c%d", 'A'+from.Col, from.Row+1)
	fmt.Println("Moving from", fromLocation)
	toLocation := fmt.Sprintf("%c%d", 'A'+to.Col, to.Row+1)
	fmt.Println("Moving to", toLocation)

	// Move the piece from the source location to the destination
	for i, piece := range b.Pieces {
		if piece.Location == fromLocation {
			b.Pieces[i].Location = toLocation
			break
		}
	}
	return false, ""
}

// getRandomPromotionPiece returns a random piece type excluding Pawn and King
func getRandomPromotionPiece() PieceType {
	promotionPieces := []PieceType{RookType, KnightType, BishopType, QueenType}
	return promotionPieces[rand.Intn(len(promotionPieces))]
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

func PositionFromString(posStr string) PossibleMovesPosition {
	col := int(posStr[0] - 'A')
	row := int(posStr[1] - '1')
	return PossibleMovesPosition{Row: row, Col: col}
}

func (b *Board) MakeRandomMove(currentColor string) (bool, string, string, PossibleMovesPosition, error) {
	validMoves = []struct {
		from PossibleMovesPosition
		to   PossibleMovesPosition
	}{}
	var kingPosition PossibleMovesPosition
	kingPositionPtr := extractKingPosition(currentColor, b)
	fmt.Println(`king position is `, kingPositionPtr)
	if kingPositionPtr != nil {
		kingPosition = *kingPositionPtr // Assign to the previously declared variable
	} else {
		return false, "", "", PossibleMovesPosition{}, fmt.Errorf("no valid moves available")
	}
	fmt.Println(`king position is `, kingPosition)

	// Collect all valid moves
	for _, piecePos := range b.Pieces {
		if piecePos.Color == currentColor {
			from := PositionFromString(piecePos.Location)
			fmt.Println(`from`, from)
			possibleMoves := b.GetPossibleMoves(from)
			for _, move := range possibleMoves {
				// Filter out moves that put the king in check
				if (move.Row == kingPosition.Row) && (move.Col == kingPosition.Col) {
					continue
				} else {
					validMoves = append(validMoves, struct {
						from PossibleMovesPosition
						to   PossibleMovesPosition
					}{from, move})
				}
			}
		}
	}

	if len(validMoves) == 0 {
		return false, "", "", PossibleMovesPosition{}, fmt.Errorf("no valid moves available")
	}

	// Randomly select a move
	fmt.Println(`validMoves`, validMoves)
	selectedMove := validMoves[rand.Intn(len(validMoves))]
	fmt.Println(`selected move`, selectedMove)

	isWin, winner, killReason := b.MovePiece(currentColor, selectedMove.from, selectedMove.to)

	return isWin, winner, killReason, selectedMove.to, nil
}

func extractKingPosition(currentColor string, b *Board) *PossibleMovesPosition {
	for _, piecePos := range b.Pieces {
		if piecePos.Type == KingType && piecePos.Color != currentColor {
			pos := PositionFromString(piecePos.Location)
			return &pos
		}
	}
	return nil
}

func extractCurrentTurnKingPosition(currentColor string, b *Board) *PossibleMovesPosition {
	for _, piecePos := range b.Pieces {
		if piecePos.Type == KingType && piecePos.Color == currentColor {
			pos := PositionFromString(piecePos.Location)
			return &pos
		}
	}
	return nil
}

func (b *Board) GetPossibleMoves(position PossibleMovesPosition) []PossibleMovesPosition {
	fmt.Println(b)
	piecePos := b.GetPiece(position)
	fmt.Println(piecePos)
	if piecePos == nil {
		return nil
	}
	// Use the pieceFactory function to create a piece instance
	piece := b.PieceFactory(piecePos.Color, piecePos.Type)
	fmt.Println(piece)
	return piece.GetPossibleMoves(b, position)
}

// PowerBoost0 converts an opponent's piece to the player's piece based on the given probabilities.
func (b *Board) PowerBoost0(currentPlayer string) error {
	fmt.Println(currentPlayer)
	// Get a list of opponent pieces on the board
	var opponentPieces []PiecePosition
	for _, piece := range b.Pieces {
		if piece.Color != currentPlayer && piece.Type != KingType {
			opponentPieces = append(opponentPieces, piece)
		}
	}

	fmt.Println("opponent pieces:", opponentPieces)

	// No opponent pieces left to convert
	if len(opponentPieces) == 0 {
		return nil
	}

	// Calculate probabilities
	totalSoldiers := 0
	totalHorses := 0
	totalBishops := 0
	totalRooks := 0
	totalQueens := 0
	for _, piece := range opponentPieces {
		switch piece.Type {
		case PawnType:
			totalSoldiers++
		case KnightType:
			totalHorses++
		case BishopType:
			totalBishops++
		case RookType:
			totalRooks++
		case QueenType:
			totalQueens++
		default:
			panic("unhandled default case")
		}
	}

	probabilities := map[PieceType]int{
		PawnType:   totalSoldiers * 25,
		KnightType: totalHorses * 25,
		BishopType: totalBishops * 25,
		RookType:   totalRooks * 15,
		QueenType:  totalQueens * 10,
	}

	// Calculate cumulative probabilities
	var cumulativeProbabilities []struct {
		pieceType PieceType
		prob      int
	}

	cumulativeSum := 0
	for pieceType, prob := range probabilities {
		cumulativeSum += prob
		cumulativeProbabilities = append(cumulativeProbabilities, struct {
			pieceType PieceType
			prob      int
		}{pieceType, cumulativeSum})
	}

	fmt.Println(currentPlayer)
	var selectedPiece *PiecePosition
	var kingPosition PossibleMovesPosition
	kingPositionPtr := extractKingPosition(currentPlayer, b)
	if kingPositionPtr != nil {
		kingPosition = *kingPositionPtr
	} else {
		return nil
	}
	fmt.Println(currentPlayer)

	// Loop to find a valid piece that doesn't threaten the king
	for {
		// Randomly select a piece type based on probabilities
		randomValue := rand.Intn(cumulativeSum)
		var selectedType PieceType
		for _, cp := range cumulativeProbabilities {
			if randomValue < cp.prob {
				selectedType = cp.pieceType
				break
			}
		}

		// Get the pieces of the selected type
		var selectedPieces []PiecePosition
		for _, piece := range opponentPieces {
			if piece.Type == selectedType {
				selectedPieces = append(selectedPieces, piece)
			}
		}

		// Randomly select one of the selected pieces
		selectedPiece = &selectedPieces[rand.Intn(len(selectedPieces))]

		from := PositionFromString(selectedPiece.Location)
		possibleMoves := b.GetPossibleMoves(from)
		valid := true

		for _, move := range possibleMoves {
			// Check if this move threatens the king
			if (move.Row == kingPosition.Row) && (move.Col == kingPosition.Col) {
				valid = false
				break
			}
		}

		if valid {
			break
		}
	}
	fmt.Println("selected piece is", selectedPiece)
	fmt.Println("currentPlayer", currentPlayer)
	// Convert the selected piece to the current player's piece
	for i, piece := range b.Pieces {
		if piece.Location == selectedPiece.Location {
			b.Pieces[i].Color = currentPlayer
			break
		}
	}

	return nil
}

func (b *Board) PowerBoost1(currentPlayer string) error {
	var playerPieces []PiecePosition
	for _, piece := range b.Pieces {
		if piece.Color == currentPlayer && piece.Type != QueenType && piece.Type != KingType {
			playerPieces = append(playerPieces, piece)
		}
	}

	// No eligible pieces left to upgrade
	if len(playerPieces) == 0 {
		return nil
	}

	// Calculate probabilities
	totalSoldiers := 0
	totalHorses := 0
	totalBishops := 0
	totalRooks := 0
	for _, piece := range playerPieces {
		switch piece.Type {
		case PawnType:
			totalSoldiers++
		case KnightType:
			totalHorses++
		case BishopType:
			totalBishops++
		case RookType:
			totalRooks++
		}
	}

	probabilities := map[PieceType]int{
		PawnType:   totalSoldiers * 25,
		KnightType: totalHorses * 25,
		BishopType: totalBishops * 25,
		RookType:   totalRooks * 15,
	}

	// Calculate cumulative probabilities
	var cumulativeProbabilities []struct {
		pieceType PieceType
		prob      int
	}

	cumulativeSum := 0
	for pieceType, prob := range probabilities {
		cumulativeSum += prob
		cumulativeProbabilities = append(cumulativeProbabilities, struct {
			pieceType PieceType
			prob      int
		}{pieceType, cumulativeSum})
	}

	var selectedPiece *PiecePosition
	var kingPosition PossibleMovesPosition
	kingPositionPtr := extractKingPosition(currentPlayer, b)
	if kingPositionPtr != nil {
		kingPosition = *kingPositionPtr
	} else {
		return nil
	}

	// Loop to find a valid piece that doesn't threaten the king
	for {
		// Randomly select a piece type based on probabilities
		randomValue := rand.Intn(cumulativeSum)
		var selectedType PieceType
		for _, cp := range cumulativeProbabilities {
			if randomValue < cp.prob {
				selectedType = cp.pieceType
				break
			}
		}

		// Get the pieces of the selected type
		var selectedPieces []PiecePosition
		for _, piece := range playerPieces {
			if piece.Type == selectedType {
				selectedPieces = append(selectedPieces, piece)
			}
		}

		// Randomly select one of the selected pieces
		selectedPiece = &selectedPieces[rand.Intn(len(selectedPieces))]

		from := PositionFromString(selectedPiece.Location)
		possibleMoves := b.GetPossibleMoves(from)
		valid := true

		for _, move := range possibleMoves {
			// Check if this move threatens the king
			if (move.Row == kingPosition.Row) && (move.Col == kingPosition.Col) {
				valid = false
				break
			}
		}

		if valid {
			break
		}
	}

	// Upgrade the selected piece
	for i, piece := range b.Pieces {
		if piece.Location == selectedPiece.Location {
			b.Pieces[i].Type = upgradePieceType(b.Pieces[i].Type)
			break
		}
	}

	return nil
}

func (b *Board) PowerBoost2(currentPlayer string) error {
	var playerPieces []PiecePosition
	for _, piece := range b.Pieces {
		if piece.Color == currentPlayer && piece.Type != KingType {
			playerPieces = append(playerPieces, piece)
		}
	}

	// Calculate probabilities for piece selection
	totalSoldiers, totalHorses, totalBishops, totalRooks, totalQueens := 0, 0, 0, 0, 0
	for _, piece := range playerPieces {
		switch piece.Type {
		case PawnType:
			totalSoldiers++
		case KnightType:
			totalHorses++
		case BishopType:
			totalBishops++
		case RookType:
			totalRooks++
		case QueenType:
			totalQueens++
		default:
			panic("unhandled default case")
		}
	}

	// Calculate cumulative probabilities
	probabilities := map[PieceType]int{
		PawnType:   totalSoldiers * 25,
		KnightType: totalHorses * 25,
		BishopType: totalBishops * 25,
		RookType:   totalRooks * 15,
		QueenType:  totalQueens * 10,
		KingType:   2, // 2% chance to teleport the king
	}

	var cumulativeProbabilities []struct {
		pieceType PieceType
		prob      int
	}

	cumulativeSum := 0
	for pieceType, prob := range probabilities {
		cumulativeSum += prob
		cumulativeProbabilities = append(cumulativeProbabilities, struct {
			pieceType PieceType
			prob      int
		}{pieceType, cumulativeSum})
	}

	for {
		randomValue := rand.Intn(cumulativeSum)
		var selectedType PieceType
		for _, cp := range cumulativeProbabilities {
			if randomValue < cp.prob {
				selectedType = cp.pieceType
				break
			}
		}

		// Get the pieces of the selected type
		var selectedPieces []PiecePosition
		for _, piece := range playerPieces {
			if piece.Type == selectedType {
				selectedPieces = append(selectedPieces, piece)
			}
		}

		// Randomly select one of the selected pieces
		selectedPiece := selectedPieces[rand.Intn(len(selectedPieces))]

		// Get all valid teleport positions
		validPositions := b.getValidTeleportPositions(selectedPiece, currentPlayer)

		if len(validPositions) == 0 {
			// No valid teleport found, try again with another piece
			continue
		}

		// Select a random valid position from the list
		teleportPosition := validPositions[rand.Intn(len(validPositions))]

		// Ensure the teleport is safe for the king
		from := PositionFromString(selectedPiece.Location)
		if selectedType == KingType {
			// For king teleport, check if the new position is safe
			if !b.isMoveSafeForKingOrTeleport(from, teleportPosition, currentPlayer, true) {
				continue // King is not safe, try another piece
			}
		} else {
			// For regular pieces, ensure the move doesn't threaten the king
			if !b.isMoveSafeForKingOrTeleport(from, teleportPosition, currentPlayer, false) {
				continue // Move threatens own king, try another piece
			}
		}

		// Apply the teleport move
		b.MovePiece(currentPlayer, from, teleportPosition)
		break // Valid move found and applied
	}
	return nil
}

// isMoveSafeForKingOrTeleport checks if moving a piece or teleporting the king keeps the king safe
func (b *Board) isMoveSafeForKingOrTeleport(from, to PossibleMovesPosition, currentPlayer string, isTeleportKing bool) bool {
	var isSafe bool // Declare isSafe outside the conditional blocks

	// Temporarily move the piece or teleport the king
	if isTeleportKing {
		// For teleporting the king, move the king to the new position
		kingPosition := extractCurrentTurnKingPosition(currentPlayer, b)
		if kingPosition == nil {
			return false // No king found
		}
		b.MovePieceForChecking(*kingPosition, to)
		isSafe = !b.isKingInCheck(currentPlayer)
		b.MovePieceForChecking(to, *kingPosition) // Revert king's move
	} else {
		// For regular piece movement
		b.MovePieceForChecking(from, to)
		isSafe = !b.isKingInCheck(currentPlayer)
		b.MovePieceForChecking(to, from) // Revert regular piece's move
	}

	return isSafe
}

// isKingInCheck checks if the king is in check after the move or teleport
func (b *Board) isKingInCheck(currentPlayer string) bool {
	kingPosition := extractCurrentTurnKingPosition(currentPlayer, b)
	if kingPosition == nil {
		return false // No king found
	}

	opponentColor := "white"
	if currentPlayer == "white" {
		opponentColor = "black"
	}

	for _, piece := range b.Pieces {
		if piece.Color == opponentColor {
			from := PositionFromString(piece.Location)
			possibleMoves := b.GetPossibleMoves(from)
			for _, move := range possibleMoves {
				if move.Row == kingPosition.Row && move.Col == kingPosition.Col {
					return true // King is in check
				}
			}
		}
	}
	return false
}

func (b *Board) getValidTeleportPositions(selectedPiece PiecePosition, currentPlayer string) []PossibleMovesPosition {
	var validPositions []PossibleMovesPosition
	from := PositionFromString(selectedPiece.Location)

	// Define forward and backward directions based on the player's color
	forwardDirection := -1 // White moves forward (rows increasing)
	backwardDirection := 1 // White moves backward (rows decreasing)
	if currentPlayer == "black" {
		forwardDirection = 1   // Black moves forward (rows decreasing)
		backwardDirection = -1 // Black moves backward (rows increasing)
	}

	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Columns; j++ {
			newPos := PossibleMovesPosition{Row: i, Col: j}

			// Determine the move direction based on a 70% chance for forward and 30% chance for backward
			direction := forwardDirection
			if rand.Float64() < 0.3 {
				direction = backwardDirection
			}

			// Skip moves that exceed one row backward
			if direction == backwardDirection && (newPos.Row-from.Row)*backwardDirection > 1 {
				continue // More than one row backward
			}

			// Ensure the piece is moving forward or staying within the valid move range
			if direction == forwardDirection && (newPos.Row-from.Row)*forwardDirection <= 0 {
				continue // Should move forward, not backward or stay in place
			}

			// Check if the position is valid and not occupied
			if b.IsValidPosition(newPos) && b.GetPiece(newPos) == nil {
				validPositions = append(validPositions, newPos)
			}
		}
	}
	return validPositions
}

func (b *Board) PlaceMine(currentPlayer string) {
	var validMinePositions []PossibleMovesPosition

	// Find the opponent's king
	var opponentKing *PiecePosition
	for _, piece := range b.Pieces {
		if piece.Color != currentPlayer && piece.Type == KingType {
			opponentKing = &piece
			break
		}
	}

	// Find valid positions for the mine, avoiding the opponent's king surrounding squares
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Columns; j++ {
			newPos := PossibleMovesPosition{Row: i, Col: j}
			location := fmt.Sprintf("%c%d", 'A'+newPos.Col, newPos.Row+1)
			if b.GetPiece(newPos) == nil { // Ensure the square is empty
				if opponentKing == nil || !isSurroundingKing(opponentKing.Location, newPos) {
					fmt.Println("possible mine position is ", location)
					validMinePositions = append(validMinePositions, newPos)
				}
			}
		}
	}

	// Select a random valid position for the mine
	if len(validMinePositions) > 0 {
		selectedPosition := validMinePositions[rand.Intn(len(validMinePositions))]
		b.Mines = append(b.Mines, Mine{Position: selectedPosition, Owner: currentPlayer, IsActive: true})
		location := fmt.Sprintf("%c%d", 'A'+selectedPosition.Col, selectedPosition.Row+1)
		fmt.Println("mine positioned at ", selectedPosition, "at location ", location)
	}
}

// Helper function to check if a position is surrounding the opponent's king
func isSurroundingKing(kingPos string, pos PossibleMovesPosition) bool {
	kingPosition := PositionFromString(kingPos)
	rowDiff := math.Abs(float64(kingPosition.Row - pos.Row))
	colDiff := math.Abs(float64(kingPosition.Col - pos.Col))
	return rowDiff <= 1 && colDiff <= 1
}

func (b *Board) TriggerMineIfPresent(to PossibleMovesPosition, currentPlayer string) (bool, string) {
	for i, mine := range b.Mines {
		// Check if the piece landed on an active mine
		if mine.Position == to && mine.IsActive && mine.Owner != currentPlayer {
			fmt.Println("Mine triggered! Destroying opponent's piece.")
			fmt.Println("To Position is")
			fmt.Println("col = ", to.Col)
			fmt.Println("row = ", to.Row)
			// Remove the opponent's piece
			for j, piece := range b.Pieces {
				if piece.Location == fmt.Sprintf("%c%d", 'A'+to.Col, to.Row+1) {
					b.Pieces = append(b.Pieces[:j], b.Pieces[j+1:]...)
					break
				}
			}

			// Deactivate the mine after triggering
			b.Mines[i].IsActive = false
			return true, "Mine triggered! Opponent's piece destroyed."
		}
	}
	// Return false if no mine was triggered
	return false, ""
}

func (b *Board) MinePowerBoost(currentPlayer string) error {
	b.PlaceMine(currentPlayer)
	return nil
}

func (b *Board) ApplyLavaStrike(currentPlayer string) ([]PossibleMovesPosition, error) {
	var lavaHitZone []PossibleMovesPosition

	// Collect all valid pieces of the current player (excluding the king)
	var playerPieces []PiecePosition
	for _, piece := range b.Pieces {
		if piece.Color == currentPlayer && piece.Type != KingType {
			playerPieces = append(playerPieces, piece)
		}
	}

	// If no valid pieces are available for the Lava strike
	if len(playerPieces) == 0 {
		return nil, fmt.Errorf("no valid pieces for Lava strike")
	}

	// Randomly select a piece for the Lava strike
	selectedPiece := playerPieces[rand.Intn(len(playerPieces))]
	selectedPosition := PositionFromString(selectedPiece.Location)

	// Determine the Lava hit zone radius (2-4 tiles)
	lavaRadius := rand.Intn(3) + 2 // Random radius between 2 and 4

	// Define possible directions with probabilities
	type Direction struct {
		offsets [][]int // Relative row and column offsets
		prob    int     // Probability weight
	}
	directions := []Direction{
		{
			offsets: [][]int{
				{1, 0}, {-1, 0}, // Up and down
			},
			prob: 30,
		},
		{
			offsets: [][]int{
				{0, 1}, {0, -1}, // Left and right
			},
			prob: 30,
		},
		{
			offsets: [][]int{
				{0, 1}, {0, -1}, {1, 0}, // Left, right, and up
			},
			prob: 30,
		},
		{
			offsets: [][]int{
				{1, 0}, {-1, 0}, {0, 1}, {0, -1}, // Up, down, left, right
			},
			prob: 10,
		},
	}

	// Calculate cumulative probabilities
	cumulativeProbabilities := []int{}
	sum := 0
	for _, dir := range directions {
		sum += dir.prob
		cumulativeProbabilities = append(cumulativeProbabilities, sum)
	}

	// Select a direction set based on the probabilities
	randomValue := rand.Intn(sum)
	var selectedDirection [][]int
	for i, prob := range cumulativeProbabilities {
		if randomValue < prob {
			selectedDirection = directions[i].offsets
			break
		}
	}

	// Generate positions within the selected direction and radius
	for _, offset := range selectedDirection {
		for step := 1; step <= lavaRadius; step++ {
			newRow := selectedPosition.Row + step*offset[0]
			newCol := selectedPosition.Col + step*offset[1]
			if b.IsValidPosition(PossibleMovesPosition{Row: newRow, Col: newCol}) {
				lavaHitZone = append(lavaHitZone, PossibleMovesPosition{Row: newRow, Col: newCol})
			}
		}
	}

	// Remove enemy pawns within the Lava hit zone
	for _, pos := range lavaHitZone {
		location := fmt.Sprintf("%c%d", 'A'+pos.Col, pos.Row+1) // Convert position to string
		for i := 0; i < len(b.Pieces); i++ {
			if b.Pieces[i].Location == location && b.Pieces[i].Color != currentPlayer && b.Pieces[i].Type == PawnType {
				// Remove the enemy pawn
				b.Pieces = append(b.Pieces[:i], b.Pieces[i+1:]...)
				break
			}
		}
	}

	// Return the Lava hit zone positions for frontend animation
	return lavaHitZone, nil
}

func (b *Board) ApplyLightningStrike(currentPlayer string) ([]PossibleMovesPosition, [][]PiecePosition, error) {
	var lightningHits []PossibleMovesPosition
	var removedPieces []PiecePosition
	var boardAfterEachHit [][]PiecePosition

	// Collect all valid target locations
	validTargets := make([]PossibleMovesPosition, 0)
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Columns; j++ {
			pos := PossibleMovesPosition{Row: i, Col: j}
			if b.IsValidPosition(pos) {
				validTargets = append(validTargets, pos)
			}
		}
	}

	if len(validTargets) == 0 {
		return nil, nil, fmt.Errorf("no valid targets for lightning strike")
	}

	// Base static probabilities
	baseProbabilities := map[PieceType]int{
		PawnType:   50,
		KnightType: 10,
		BishopType: 15,
		RookType:   8,
		QueenType:  7,
		KingType:   0, // King cannot be killed
	}

	emptySquareProbability := 10 // Default probability for empty squares

	// Adjust probabilities dynamically based on available pieces
	totalPieces := len(b.Pieces)
	if totalPieces <= 5 {
		// Increase odds for empty squares when fewer pieces remain
		emptySquareProbability = 70
	}

	adjustedProbabilities := calculateAdjustedProbabilities(baseProbabilities, b.Pieces, currentPlayer)

	minKillReached := false
	playerKills := 0
	strikes := 3

	for strikes > 0 {
		// Select a random target position
		targetIdx := rand.Intn(len(validTargets))
		targetPos := validTargets[targetIdx]
		validTargets = append(validTargets[:targetIdx], validTargets[targetIdx+1:]...) // Remove target from the list

		// Determine if the target is a piece or an empty square
		piece := b.GetPiece(targetPos)

		if piece != nil {
			// Get the probability of killing the piece
			pieceOdds := adjustedProbabilities[piece.Type]
			if piece.Color == currentPlayer {
				pieceOdds = pieceOdds / 3 // Reduce odds for killing player's own pieces
			}

			if rand.Intn(100) < pieceOdds {
				// Kill the piece
				if piece.Color != currentPlayer || playerKills < 1 {
					removedPieces = append(removedPieces, *piece)
					b.RemovePiece(targetPos) // Remove piece from the board
					if piece.Color == currentPlayer {
						playerKills++
					}
					minKillReached = true
				}
			}
		} else if rand.Intn(100) < emptySquareProbability {
			// Hit an empty square (no action needed for now)
		}

		lightningHits = append(lightningHits, targetPos)
		boardAfterEachHit = append(boardAfterEachHit, b.Pieces)
		strikes--
		fmt.Println(`lightning is `, lightningHits)
		// If no valid targets remain, break the loop
		if len(validTargets) == 0 {
			break
		}
	}

	// If no kills occurred after all strikes, repeat if valid targets remain
	if !minKillReached {
		if len(validTargets) == 0 {
			return nil, nil, fmt.Errorf("no valid targets remaining for lightning strike")
		}
		return b.ApplyLightningStrike(currentPlayer)
	}

	return lightningHits, boardAfterEachHit, nil
}

// Helper function to calculate adjusted probabilities based on available pieces
func calculateAdjustedProbabilities(baseProbabilities map[PieceType]int, pieces []PiecePosition, currentPlayer string) map[PieceType]int {
	adjustedProbabilities := make(map[PieceType]int)
	availablePieceCounts := make(map[PieceType]int)

	// Count available pieces of each type
	for _, piece := range pieces {
		if piece.Color != currentPlayer || piece.Type == KingType {
			continue // Exclude the player's pieces and the King
		}
		availablePieceCounts[piece.Type]++
	}

	// Adjust probabilities based on available pieces
	totalProbability := 0
	for pieceType, baseProb := range baseProbabilities {
		if availablePieceCounts[pieceType] == 0 {
			// Distribute the probability of missing pieces equally among others
			continue
		} else {
			adjustedProbabilities[pieceType] = baseProb
			totalProbability += baseProb
		}
	}

	// Normalize probabilities to ensure they sum to 100
	for pieceType := range adjustedProbabilities {
		adjustedProbabilities[pieceType] = (adjustedProbabilities[pieceType] * 100) / totalProbability
	}

	return adjustedProbabilities
}

// RemovePiece Helper method to remove a piece from the board
func (b *Board) RemovePiece(pos PossibleMovesPosition) {
	location := fmt.Sprintf("%c%d", 'A'+pos.Col, pos.Row+1)
	for i, piece := range b.Pieces {
		if piece.Location == location {
			b.Pieces = append(b.Pieces[:i], b.Pieces[i+1:]...)
			break
		}
	}
}

func upgradePieceType(currentType PieceType) PieceType {
	switch currentType {
	case PawnType:
		return KnightType
	case KnightType:
		return BishopType
	case BishopType:
		return RookType
	case RookType:
		return QueenType
	default:
		return currentType
	}
}
