package engine

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func PlayAI(rw io.ReadWriter) {
	reader := bufio.NewReader(rw)

	initBoard := Board{
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
	}

	gameConfig := &Config{Board: initBoard}

	readLine := func() (string, bool) {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", false
		}
		return strings.TrimSpace(line), true
	}

	userPlayerInput := ""
	for {
		fmt.Fprint(rw, "Are you X or O? > ")
		line, ok := readLine()
		if !ok {
			return
		}
		if line != "X" && line != "O" {
			fmt.Fprintln(rw, "enter a valid choice!")
			continue
		}
		userPlayerInput = line
		break
	}

	userPlayer := stringToMove(userPlayerInput)
	fmt.Fprintf(rw, "You are %s. Let's start the game!\n", userPlayerInput)

	for {
		fmt.Fprint(rw, "Player input >> ")

		userInput, ok := readLine()
		if !ok {
			return
		}

		win, err := PlayMove(userInput, gameConfig, userPlayer, "Player")
		if err != nil {
			fmt.Fprintf(rw, "error: %v\n", err)
			continue
		}
		PrintBoard(rw, gameConfig.Board)

		if win != "" {
			if win == "Draw" {
				fmt.Fprintln(rw, "The Game has ended as a draw!")
				return
			}

			fmt.Fprintf(rw, "%s wins!\n", win)
			return
		}

		aiPlayer := EMPTY
		if userPlayer == X {
			aiPlayer = O
		} else {
			aiPlayer = X
		}

		fmt.Fprintln(rw, "AI plays!")

		action := Minimax(gameConfig.Board, userPlayer)
		aiMove := coordToInt(action)

		win, _ = PlayMove(strconv.Itoa(aiMove), gameConfig, aiPlayer, "AI")
		PrintBoard(rw, gameConfig.Board)

		if win != "" {
			if win == "Draw" {
				fmt.Fprintln(rw, "The Game has ended as a draw!")
				return
			}

			fmt.Fprintf(rw, "%s wins!\n", win)
			return
		}
	}
}
