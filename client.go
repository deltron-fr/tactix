package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	"github.com/deltron-fr/tactix/internal/engine"
	"github.com/gorilla/websocket"
)

func createClient() *websocket.Conn {
	URL := "ws://localhost:8080/ws"
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(URL, nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}

	return conn

}

func player(scanner *bufio.Scanner) {
	conn := createClient()
	defer conn.Close()

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("\nServer >> %s\nPlayer input >> ", msg)
		}
	}()

	for {
		fmt.Printf("Player input >> ")

		if scanner.Scan() {
			userInput := scanner.Text()
			err := conn.WriteMessage(websocket.TextMessage, []byte(userInput))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func playAI(scanner *bufio.Scanner) {

	initBoard := engine.Board{
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
	}

	gameConfig := &engine.Config{Board: initBoard}

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
