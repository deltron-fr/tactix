package main

import (
	"errors"
	"strconv"
)

type Board [][]move

type Config struct {
	board Board
}

func playMove(userMove string, cfg *Config, gamePlayer move, playerName string) (string, error) {

	// char, pos := r[0], int(r[1] - '0')

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

		cfg.board[0][0] = gamePlayer
	case 2:
		err := verifyMove(cfg, 0, 1)
		if err != nil {
			return "", err
		}
		
		cfg.board[0][1] = gamePlayer
	case 3:
		err := verifyMove(cfg, 0, 2)
		if err != nil {
			return "", err
		}
		
		cfg.board[0][2] = gamePlayer
	case 4:
		err := verifyMove(cfg, 1, 0)
		if err != nil {
			return "", err
		}
		
		cfg.board[1][0] = gamePlayer
	case 5:
		err := verifyMove(cfg, 1, 1)
		if err != nil {
			return "", err
		}
		
		cfg.board[1][1] = gamePlayer
	case 6:
		err := verifyMove(cfg, 1, 2)
		if err != nil {
			return "", err
		}
		
		cfg.board[1][2] = gamePlayer
	case 7:
		err := verifyMove(cfg, 2, 0)
		if err != nil {
			return "", err
		}
		
		cfg.board[2][0] = gamePlayer
	case 8:
		err := verifyMove(cfg, 2, 1)
		if err != nil {
			return "", err
		}
		
		cfg.board[2][1] = gamePlayer
	case 9:
		err := verifyMove(cfg, 2, 2)
		if err != nil {
			return "", err
		}
		
		cfg.board[2][2] = gamePlayer
	}

	gameWinner := ""
	if terminal(cfg.board) {
		winner := winner(cfg.board)
		switch winner {
		case EMPTY:
			gameWinner = "Draw"
		default:
			gameWinner = playerName
		}	
	}

	return gameWinner, nil
}

func verifyMove(cfg *Config, row, col int) error {
	if cfg.board[row][col] != EMPTY {
		return errors.New("invalid move")
	}

	return nil
}