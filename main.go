package main

import "fmt"

func main() {
	fmt.Println("Chess Engine v1.0")
	fmt.Println("===================")

	board := NewBoard()
	board.Display()

	// Generate Moves for White
	whiteMoves := board.GenerateAllMoves(White)
	fmt.Printf("\nWhite has %d possible moves:\n", len(whiteMoves))

	// Show first 10 moves
	for i, move := range whiteMoves {
		if i > 10 {
			fmt.Println("...")
			break
		}
		fmt.Printf("%d. %s\n", i+1, move.String())
	}

	// Make move and print result
	fmt.Println("\nAfter e2e4:")
	move := Move{FromRow: 6, FromCol: 4, ToRow: 4, ToCol: 4, PieceType: Pawn}
	board.MakeMove(move)
	board.Display()

}
