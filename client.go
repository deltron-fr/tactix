package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	"github.com/deltron-fr/tactix/internal/engine"
	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	game       *Game
}

func NewClient(conn *websocket.Conn, manager *Manager, game *Game) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		game:       game,
	}
}

func createClient() {
	URL := "ws://localhost:8080/ws"
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(URL, nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer conn.Close()

}

func player(scanner *bufio.Scanner) {
	URL := "ws://localhost:8080/ws"
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(URL, nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer conn.Close()

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

		go func() {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("\nServer >> %s\nPlayer input >> ", msg)
		}()
	}

	
}

func (c *Client) readMessages() {
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close: %v", err)

			}
			break
		}
		log.Println(messageType)
		log.Println(payload)
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
