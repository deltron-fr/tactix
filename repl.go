package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deltron-fr/tactix/internal/engine"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("================ Welcome to TacTix! ===============")
	fmt.Println("=================================================")
	fmt.Println("		1  |  2  |  3		")
	fmt.Println("		-------------		")
	fmt.Println("		4  |  5  |  6		")
	fmt.Println("		--------------		")
	fmt.Println("		7  |  8  |  9		")

	fmt.Println("Game has started!")

	fmt.Println("Choose a mode: ")
	fmt.Println("1. Play against another player")
	fmt.Println("2. Play vs AI")

	mode := ""
	for {
		fmt.Printf("1 or 2 > ")
		if scanner.Scan() {
			mode = scanner.Text()
			if mode != "1" && mode != "2" {
				fmt.Println("enter a valid choice!")
			} else {
				break
			}
		}
	}

	if mode == "2" {
		playAI(scanner)
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
