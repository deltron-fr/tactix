package main

import (
	"slices"
)

type move int

const (
	EMPTY move = iota
	X 
	O
)

func (m move) String() string {
    switch m {
    case X:
        return "X"
    case O:
        return "O"
    default:
        return "-" 
    }
}


func winner(board Board) move {
	winningLines := [][][]int{
		{{0,0},{0,1},{0,2}}, // row 0
		{{1,0},{1,1},{1,2}}, // row 1
		{{2,0},{2,1},{2,2}}, // row 2
		{{0,0},{1,0},{2,0}}, // col 0
		{{0,1},{1,1},{2,1}}, // col 1
		{{0,2},{1,2},{2,2}}, // col 2
		{{0,0},{1,1},{2,2}}, // diag 1
		{{0,2},{1,1},{2,0}}, // diag 2
	}

	for _, line := range winningLines {
		first := board[line[0][0]][line[0][1]]

		if first == EMPTY {
			continue
		}

		win := true

		for _, l := range line {
			if board[l[0]][l[1]] != first {
				win = false
				break
			}
		}

		if win {
			return first
		}

	}

	return EMPTY
}

func terminal(board Board) bool {
	
	if winner(board) != EMPTY {
		return true
	}
	
	emptyCell := false

	for _, row := range board {
		if slices.Contains(row, EMPTY) {
			emptyCell = true
		}
	}
	
	return !emptyCell
}



