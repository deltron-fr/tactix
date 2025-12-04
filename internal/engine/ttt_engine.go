package engine

import (
	"slices"
	"strconv"
	"errors"
	"math"
)


func Minimax(board Board, initMove Move) []int {
	// Returns the optimal action for the current player on the board.

	if player(board, initMove) == X {
		_, bestMove := maxValue(board, initMove)
		return bestMove
	}

	_, bestMove := minValue(board, initMove)
	return bestMove
}

func PlayMove(userMove string, cfg *Config, gamePlayer Move, playerName string) (string, error) {
	// Verifies the players move, saves the player(or AI)s move to the board and returns a game winner if any.

	pos, err := strconv.Atoi(userMove)
	if err != nil {
		return "", errors.New("input a valid number")
	}


	if pos < 1 || pos > 9 {
		return "", errors.New("number isn't a valid position on the board")
	}

	switch pos{
	case 1:
		err := verifyMove(cfg, 0, 0)
		if err != nil {
			return "", err
		}

		cfg.Board[0][0] = gamePlayer
	case 2:
		err := verifyMove(cfg, 0, 1)
		if err != nil {
			return "", err
		}
		
		cfg.Board[0][1] = gamePlayer
	case 3:
		err := verifyMove(cfg, 0, 2)
		if err != nil {
			return "", err
		}
		
		cfg.Board[0][2] = gamePlayer
	case 4:
		err := verifyMove(cfg, 1, 0)
		if err != nil {
			return "", err
		}
		
		cfg.Board[1][0] = gamePlayer
	case 5:
		err := verifyMove(cfg, 1, 1)
		if err != nil {
			return "", err
		}
		
		cfg.Board[1][1] = gamePlayer
	case 6:
		err := verifyMove(cfg, 1, 2)
		if err != nil {
			return "", err
		}
		
		cfg.Board[1][2] = gamePlayer
	case 7:
		err := verifyMove(cfg, 2, 0)
		if err != nil {
			return "", err
		}
		
		cfg.Board[2][0] = gamePlayer
	case 8:
		err := verifyMove(cfg, 2, 1)
		if err != nil {
			return "", err
		}
		
		cfg.Board[2][1] = gamePlayer
	case 9:
		err := verifyMove(cfg, 2, 2)
		if err != nil {
			return "", err
		}
		
		cfg.Board[2][2] = gamePlayer
	}

	gameWinner := ""
	if terminal(cfg.Board) {
		winner := winner(cfg.Board)
		switch winner {
		case EMPTY:
			gameWinner = "Draw"
		default:
			gameWinner = playerName
		}	
	}

	return gameWinner, nil
}



func player(board Board, initPlayer Move) Move {
	// Returns player who has the next turn on a board.

	xMoves := 0
	oMoves := 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch board[i][j] {
			case X:
				xMoves++
			case O:
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
	// Returns set of all possible actions (i, j) available on the board.

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


func maxValue(board Board, initMove Move) (int, []int) {
	// Recursively gets the max value from all possible actions for a given board state,
	//  then returns the best action(max value).
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

func minValue(board Board, initMove Move) (int, []int) {
	// Recursively gets the min value from all possible actions for a given board state,
	//  then returns the best action(min value).
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


func result(board Board, action []int, initMove Move) (Board, error) {
	// Returns the board that results from making move (i, j) on the board.
	boardCopy := make(Board, len(board))
	for i := range board {
		rowCopy := make([]Move, len(board[i]))
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


func terminal(board Board) bool {
	// Returns true if game is over, false if it isn't.
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

func utility(board Board) int {
	// Returns the numerical value for final board state, 
	// X wins - 1, O wins - (-1) and a draw - 0
	if winner(board) == X {
		return 1
	} else if winner(board) == O {
		return -1
	}

	return 0
}

func winner(board Board) Move {
	// Returns winner of the gamme if any
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


func verifyMove(cfg *Config, row, col int) error {
	// Utility function to verify the players move
	if cfg.Board[row][col] != EMPTY {
		return errors.New("invalid move")
	}

	return nil
}
