package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)


func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)

	initBoard := Board{
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
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

	userPlayerInput := ""

	for {
		fmt.Printf("Are you X or O? > ")
		if scanner.Scan() {
			userPlayerInput = scanner.Text()
			if userPlayerInput != "X" && userPlayerInput != "O" {
				fmt.Println("enter a valid choice!")
			} else {
				break
			}
		}
	}
	
	userPlayer := stringToMove(userPlayerInput)

	for {	
		fmt.Printf("Player input >> ")

		if scanner.Scan() {
			userInput := scanner.Text()
			win, err := playMove(userInput, gameConfig, userPlayer, "Player")
			if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			printBoard(gameConfig)

			if win != "" {
				if win == "Draw" {
					fmt.Println("The Game has ended as a draw!")
					return
				}

				fmt.Printf("%s wins!\n", win)
				return
			}

		}

		aiPlayer := EMPTY
		if userPlayer == X {
			aiPlayer = O
		} else {
			aiPlayer = X
		}

		fmt.Printf("AI plays!\n")

		action := minimax(gameConfig.board, userPlayer)
		aiMove := coordToInt(action)
		
		win, _ := playMove(strconv.Itoa(aiMove), gameConfig, aiPlayer, "AI")
		printBoard(gameConfig)

		if win != "" {
				if win == "Draw" {
					fmt.Println("The Game has ended as a draw!")
					return
				}

				fmt.Printf("%s wins!\n", win)
				return
			}
		}
}


func printBoard(cfg *Config) {

	for i := 0; i < 3; i++ {
		for i := 0; i < 3; i++ {
			if cfg.board[i][i] != X && cfg.board[i][i] != O {
						cfg.board[i][i] = EMPTY
					}
		}
	}


	fmt.Printf("		%s  |  %s  |  %s		\n", cfg.board[0][0].String(), cfg.board[0][1].String(), cfg.board[0][2].String())
	fmt.Println("		-------------")
	fmt.Printf("		%s  |  %s  |  %s		\n", cfg.board[1][0].String(), cfg.board[1][1].String(), cfg.board[1][2].String())
	fmt.Println("		--------------	")
	fmt.Printf("		%s  |  %s  |  %s		\n", cfg.board[2][0].String(), cfg.board[2][1].String(), cfg.board[2][2].String())
	fmt.Println()
}

func stringToMove(input string) move {
	switch input {
	case X.String():
		return X
	case O.String():
		return O
	}

	return EMPTY
}