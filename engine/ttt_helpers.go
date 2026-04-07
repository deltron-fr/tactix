package engine

import (
	"fmt"
	"io"
)

func coordToInt(action []int) int {
	value := 0

	switch action[0] {
	case 0:
		if action[1] == 0 {
			value = 1
		} else if action[1] == 1 {
			value = 2
		} else if action[1] == 2 {
			value = 3
		}
	case 1:
		if action[1] == 0 {
			value = 4
		} else if action[1] == 1 {
			value = 5
		} else if action[1] == 2 {
			value = 6
		}
	case 2:
		if action[1] == 0 {
			value = 7
		} else if action[1] == 1 {
			value = 8
		} else if action[1] == 2 {
			value = 9
		}
	}

	return value
}

func PrintBoard(w io.Writer, cfg *Config) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if cfg.Board[i][j] != X && cfg.Board[i][j] != O {
				cfg.Board[i][j] = EMPTY
			}
		}
	}

	fmt.Fprintf(w, "		%s  |  %s  |  %s		\n", cfg.Board[0][0].String(), cfg.Board[0][1].String(), cfg.Board[0][2].String())
	fmt.Fprintln(w, "		-------------")
	fmt.Fprintf(w, "		%s  |  %s  |  %s		\n", cfg.Board[1][0].String(), cfg.Board[1][1].String(), cfg.Board[1][2].String())
	fmt.Fprintln(w, "		--------------	")
	fmt.Fprintf(w, "		%s  |  %s  |  %s		\n", cfg.Board[2][0].String(), cfg.Board[2][1].String(), cfg.Board[2][2].String())
	fmt.Fprintln(w)
}

func stringToMove(input string) Move {
	switch input {
	case X.String():
		return X
	case O.String():
		return O
	}

	return EMPTY
}

