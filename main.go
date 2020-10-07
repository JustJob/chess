package main

import "context"

func main() {
	g := NewStandardChessGame()
	g.PlayGame(context.Background())
}
