package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/deltron-fr/tactix/internal/engine"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	game       *Game
	symbol     engine.Move
}

func NewClient(conn *websocket.Conn, manager *Manager, game *Game, symbol engine.Move) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		game:       game,
		symbol:     symbol,
	}
}

type clientQueue struct {
	Elements []*Client
}

func (c clientQueue) IsEmpty() bool {
	return len(c.Elements) == 0
}

func (c clientQueue) Size() int {
	return len(c.Elements)
}

func (c *clientQueue) Dequeue() (*Client, bool) {
	if c.IsEmpty() {
		return nil, false
	}

	ele := c.Elements[0]
	c.Elements = c.Elements[1:]

	return ele, true
}

func (c *clientQueue) Enqueue(item *Client) {
	c.Elements = append(c.Elements, item)
}

type Manager struct {
	queue *clientQueue
	mu    sync.Mutex
}

type Message struct {
	Sender  *Client
	Payload []byte
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	switch m.queue.Size() {
	case 0:
		client := NewClient(conn, m, nil, engine.X)
		m.addClient(client, m.queue)
	case 1:
		client := NewClient(conn, m, nil, engine.O)
		m.addClient(client, m.queue)
	}

	m.mu.Lock()
	if m.queue.Size() == 2 {
		players := make(map[*Client]bool)
		firstClient, _ := m.queue.Dequeue()
		secondClient, _ := m.queue.Dequeue()
		players[firstClient] = true
		players[secondClient] = true

		initBoard := engine.Board{
			{engine.EMPTY, engine.EMPTY, engine.EMPTY},
			{engine.EMPTY, engine.EMPTY, engine.EMPTY},
			{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		}

		gameConfig := &engine.Config{
			Board:   initBoard,
			Player1: engine.X,
			Player2: engine.O,
		}

		newGame := Game{
			Clients:   players,
			GameState: gameConfig,
		}

		firstClient.game = &newGame
		secondClient.game = &newGame

		// go m.handleMessages(players, broadcast, &newGame)

	}
	m.mu.Unlock()
}

func (c *Client) handleClients() {
	msgChan := make(chan Message)
	for {
		_, msg, err := c.connection.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		clientMessage := Message{
			Sender:  c,
			Payload: msg,
		}
		msgChan <- clientMessage
	}
}

func (m *Manager) handleMessages(players map[*Client]bool, msg chan []byte, gameState *Game) {
	for {
		message := <-msg
		// validateMove(string(message), gameState.GameState)

		for c, _ := range players {
			err := c.connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				return
			}

		}
	}
}

func validateMove(userMove string, cfg *engine.Config, gamePlayer engine.Move, playerName string) (string, error) {

	pos, err := strconv.Atoi(userMove)
	if err != nil {
		return "", errors.New("input a valid number")
	}

	if pos < 1 || pos > 9 {
		return "", errors.New("number isn't a valid position on the board")
	}

	switch pos {
	case 1:
		err := engine.VerifyMove(cfg, 0, 0)
		if err != nil {
			return "", err
		}

		cfg.Board[0][0] = gamePlayer
	case 2:
		err := engine.VerifyMove(cfg, 0, 1)
		if err != nil {
			return "", err
		}

		cfg.Board[0][1] = gamePlayer
	case 3:
		err := engine.VerifyMove(cfg, 0, 2)
		if err != nil {
			return "", err
		}

		cfg.Board[0][2] = gamePlayer
	case 4:
		err := engine.VerifyMove(cfg, 1, 0)
		if err != nil {
			return "", err
		}

		cfg.Board[1][0] = gamePlayer
	case 5:
		err := engine.VerifyMove(cfg, 1, 1)
		if err != nil {
			return "", err
		}

		cfg.Board[1][1] = gamePlayer
	case 6:
		err := engine.VerifyMove(cfg, 1, 2)
		if err != nil {
			return "", err
		}

		cfg.Board[1][2] = gamePlayer
	case 7:
		err := engine.VerifyMove(cfg, 2, 0)
		if err != nil {
			return "", err
		}

		cfg.Board[2][0] = gamePlayer
	case 8:
		err := engine.VerifyMove(cfg, 2, 1)
		if err != nil {
			return "", err
		}

		cfg.Board[2][1] = gamePlayer
	case 9:
		err := engine.VerifyMove(cfg, 2, 2)
		if err != nil {
			return "", err
		}

		cfg.Board[2][2] = gamePlayer
	}

	gameWinner := ""
	if engine.Terminal(cfg.Board) {
		winner := engine.Winner(cfg.Board)
		switch winner {
		case engine.EMPTY:
			gameWinner = "Draw"
		default:
			gameWinner = playerName
		}
	}

	return gameWinner, nil
}

func (m *Manager) addClient(client *Client, q *clientQueue) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q.Enqueue(client)
}
