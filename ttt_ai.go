package main

import (
	"errors"
	"math"
	"slices"
)

func player(board Board, initPlayer move) move {
	xMoves := 0
	oMoves := 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == X {
				xMoves++
			} else if board[i][j] == O {
				oMoves++
			}
		}
	}

	if xMoves > oMoves {
		return O
	} else if oMoves > xMoves {
		return X
	}

	return initPlayer
}

func actions(board Board) [][]int {
	var actions [][]int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] != X && board[i][j] != O {
				pos := []int{i, j}
				actions = append(actions, pos)
			}
		}
	}

	return actions
}

func result(board Board, action []int, initMove move) (Board, error) {

	boardCopy := make(Board, len(board))
	for i := range board {
		rowCopy := make([]move, len(board[i]))
		copy(rowCopy, board[i])
		boardCopy[i] = rowCopy
	}

	found := false

	actions := actions(boardCopy)
	for i := 0; i < len(actions); i++ {
		if slices.Equal(action, actions[i]) {
			found = true
		}
	}

	if !found {
		return boardCopy, errors.New("action is not possible")
	}

	boardCopy[action[0]][action[1]] = player(boardCopy, initMove)

	return boardCopy, nil
}


func utility(board Board) int {
	if winner(board) == X {
		return 1
	} else if winner(board) == O {
		return -1
	}

	return 0
}

func minimax(board Board, initMove move) []int {

	if player(board, initMove) == X {
		_, bestMove := maxValue(board, initMove)
		return bestMove
	}

	_, bestMove := minValue(board, initMove)
	return bestMove
}

func maxValue(board Board, initMove move) (int, []int) {
	optimalAction := []int{}

	if terminal(board) {
		return utility(board), optimalAction
	}

	value := math.MinInt64

	for _, action := range(actions(board)) {
		currentValue := value

		newBoard, err := result(board, action, initMove)
		if err != nil {
			return 0, nil
		}

		minV, _ := minValue(newBoard, initMove)
		value = max(value, minV)
		if currentValue != value {
			optimalAction = action
		}
	}

	return value, optimalAction

}

func minValue(board Board, initMove move) (int, []int) {
	optimalAction := []int{}

	if terminal(board) {
		return utility(board), optimalAction
	}

	value := math.MaxInt64

	for _, action := range(actions(board)) {
		currentValue := value

		newBoard, err := result(board, action, initMove)
		if err != nil {
			return 0, nil
		}

		maxV, _ := maxValue(newBoard, initMove)
		value = min(value, maxV)

		if currentValue != value {
			optimalAction = action
		}
	}

	return value, optimalAction
}

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