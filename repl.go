package main

import (
	"bufio"
	"fmt"
	"os"
)


func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)

	initBoard := Board{
		{'-', '-', '-'},
		{'-', '-', '-'},
		{'-', '-', '-'},
	}
	
	gameConfig := &Config{board: initBoard}

	fmt.Println("================ Welcome to TacTix! ===============")
	fmt.Println("=================================================")
	fmt.Println("		1  |  2  |  3		")
	fmt.Println("		-------------		")
	fmt.Println("		4  |  5  |  6		")
	fmt.Println("		--------------		")
	fmt.Println("		7  |  8  |  9		")

	fmt.Println("Game has started!")
	for {
		
		fmt.Printf("player 1 input >> ")

		if scanner.Scan() {
			userInput := scanner.Text()
			_ = playMove(userInput, gameConfig)

			printBoard(gameConfig)
			
			}
		fmt.Printf("player 2 input >> ")

		if scanner.Scan() {
			userInput := scanner.Text()
			_ = playMove(userInput, gameConfig)

			printBoard(gameConfig)

		}
	}
}


func printBoard(cfg *Config) {

	for i := 0; i < 3; i++ {
		for i := 0; i < 3; i++ {
			if cfg.board[i][i] != 'X' && cfg.board[i][i] != 'O' {
						cfg.board[i][i] = '-'
					}
		}
	}


	fmt.Printf("		%c  |  %c  |  %c		\n", cfg.board[0][0], cfg.board[0][1], cfg.board[0][2])
	fmt.Println("		-------------")
	fmt.Printf("		%c  |  %c  |  %c		\n", cfg.board[1][0], cfg.board[1][1], cfg.board[1][2])
	fmt.Println("		--------------	")
	fmt.Printf("		%c  |  %c  |  %c		\n", cfg.board[2][0], cfg.board[2][1], cfg.board[2][2])
	fmt.Println()
}