package main

type Board [][]rune

type Config struct {
	board Board
}

func playMove(userMove string, cfg *Config) *Board {
	r := []rune(userMove)

	char, pos := r[0], int(r[1] - '0')

	switch pos{
	case 1:
		cfg.board[0][0] = char
	case 2:
		cfg.board[0][1] = char
	case 3:
		cfg.board[0][2] = char
	case 4:
		cfg.board[1][0] = char
	case 5:
		cfg.board[1][1] = char
	case 6:
		cfg.board[1][2] = char
	case 7:
		cfg.board[2][0] = char
	case 8:
		cfg.board[2][1] = char
	case 9:
		cfg.board[2][2] = char
	}

	return &cfg.board
}