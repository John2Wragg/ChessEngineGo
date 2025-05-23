package main

import "fmt"

func main() {
	fmt.Println("Chess Engine v1.0")
	fmt.Println("===================")

	board := NewBoard()
	board.Display()

	// Generate Moves for White
	whiteMoves := board.GenerateLegalMoves(White)
	fmt.Printf("\nWhite has %d legal moves:\n", len(whiteMoves))

	// Show first 10 moves
	for i, move := range whiteMoves {
		if i > 10 {
			fmt.Println("...")
			break
		}
		fmt.Printf("%d. %s\n", i+1, move.String())
	}

	// Test game state
	fmt.Printf("\nGame State: %s\n", board.GetGameResult(White))

	// Make move and print result
	fmt.Println("\nAfter e2e4, e7e5:")
	board.MakeMove(Move{FromRow: 6, FromCol: 4, ToRow: 4, ToCol: 4, PieceType: Pawn})
	board.MakeMove(Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4, PieceType: Pawn})
	board.Display()

	whiteMoves = board.GenerateLegalMoves(White)
	fmt.Printf("\nWhite now has %d legal moves", len(whiteMoves))

}
