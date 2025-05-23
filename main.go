package main

import "fmt"

func main() {
	fmt.Println("Chess Engine v1.0")
	fmt.Println("===================")

	board := NewBoard()
	board.Display()
}
