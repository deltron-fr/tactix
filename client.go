package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/deltron-fr/tactix/internal/engine"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	initBoard := engine.Board{
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
	}

	gameConfig := &engine.Config{Board: initBoard}

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
			win, err := engine.PlayMove(userInput, gameConfig, userPlayer, "Player")
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

		aiPlayer := engine.EMPTY
		if userPlayer == engine.X {
			aiPlayer = engine.O
		} else {
			aiPlayer = engine.X
		}

		fmt.Printf("AI plays!\n")

		action := engine.Minimax(gameConfig.Board, userPlayer)
		aiMove := coordToInt(action)

		win, _ := engine.PlayMove(strconv.Itoa(aiMove), gameConfig, aiPlayer, "AI")
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

func printBoard(cfg *engine.Config) {

	for i := 0; i < 3; i++ {
		for i := 0; i < 3; i++ {
			if cfg.Board[i][i] != engine.X && cfg.Board[i][i] != engine.O {
				cfg.Board[i][i] = engine.EMPTY
			}
		}
	}

	fmt.Printf("		%s  |  %s  |  %s		\n", cfg.Board[0][0].String(), cfg.Board[0][1].String(), cfg.Board[0][2].String())
	fmt.Println("		-------------")
	fmt.Printf("		%s  |  %s  |  %s		\n", cfg.Board[1][0].String(), cfg.Board[1][1].String(), cfg.Board[1][2].String())
	fmt.Println("		--------------	")
	fmt.Printf("		%s  |  %s  |  %s		\n", cfg.Board[2][0].String(), cfg.Board[2][1].String(), cfg.Board[2][2].String())
	fmt.Println()
}

func stringToMove(input string) engine.Move {
	switch input {
	case engine.X.String():
		return engine.X
	case engine.O.String():
		return engine.O
	}

	return engine.EMPTY
}
