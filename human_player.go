package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

type HumanPlayer struct {
}

func (h *HumanPlayer) GetNextMove(ctx context.Context, board *Board) (*Move, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("What is your next move: ")
	if scanner.Scan() {
		mv := scanner.Text()
		m, err := MoveFromStringPos(mv)
		if err != nil {
			return nil, err
		}

		return m, nil
	}

	return nil, errors.New("scan failed")
}

func (h *HumanPlayer) GetTimePerTurn() time.Duration {
	return 30 * time.Second
}
