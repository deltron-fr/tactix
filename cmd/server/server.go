package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/deltron-fr/tactix/client"
	"github.com/deltron-fr/tactix/engine"
)

const (
	port    = ":8000"
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type GameInfo struct {
	Clients map[client.Client]bool
	Mu      sync.Mutex
}

func newGameInfo() *GameInfo {
	return &GameInfo{
		Clients: make(map[client.Client]bool),
		Mu:      sync.Mutex{},
	}
}

type Room struct {
	ID        string
	GameBoard engine.Board
	Players   [2]client.Client
	Turn      int
	Manager   *Manager
	Mu        sync.Mutex
}

func generateRoomID() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return "ROOM-" + string(b)
}

func NewRoom(manager *Manager, player1, player2 client.Client) *Room {
	initBoard := engine.Board{
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
		{engine.EMPTY, engine.EMPTY, engine.EMPTY},
	}

	return &Room{
		ID:        generateRoomID(),
		GameBoard: initBoard,
		Players:   [2]client.Client{player1, player2},
		Turn:      0,
		Manager:   manager,
		Mu:        sync.Mutex{},
	}
}

var matchMaker = make(chan client.Client, 2)

func startServer() {
	l, err := net.Listen("tcp", "localhost"+port)
	if err != nil {
		log.Fatal(err)
	}

	gameInfo := newGameInfo()
	mgr := NewManager()
	go mgr.handleRoomCleanup()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("ERROR - Error accepting connection:", err)
			continue
		}

		cli := client.Client(conn)

		gameInfo.Mu.Lock()
		gameInfo.Clients[cli] = true
		gameInfo.Mu.Unlock()

		mgr.handleMatchMaking(cli)
	}
}

func (m *Manager) handleMatchMaking(cli client.Client) {
	fmt.Printf("New client connected: %v\n", cli.RemoteAddr())
	cli.Write([]byte("Welcome to TacTix! Waiting for an opponent...\n"))
	matchMaker <- cli

	if len(matchMaker) < 2 {
		return
	}

	var player1, player2 client.Client

	player1 = <-matchMaker
	player2 = <-matchMaker

	room := NewRoom(m, player1, player2)
	go handleGame(room)
}

func send(c client.Client, format string, a ...any) {
	fmt.Fprintf(c, format, a...)
}

func handleGame(room *Room) {
	readers := [2]*bufio.Reader{
		bufio.NewReader(room.Players[0]),
		bufio.NewReader(room.Players[1]),
	}
	marks := [2]engine.Move{engine.X, engine.O}

	send(room.Players[0], "Game starting in room %s. You are X.\n", room.ID)
	send(room.Players[1], "Game starting in room %s. You are O.\n", room.ID)

	for {
		idx := room.Turn
		player := room.Players[idx]
		other := room.Players[1-idx]

		send(player, "Your turn (%s). Enter a position 1-9:\n", marks[idx])
		send(other, "Waiting for opponent...\n")

		line, err := readers[idx].ReadString('\n')
		if err != nil {
			log.Printf("read error from player %d: %v", idx+1, err)
			send(other, "Opponent disconnected. Game over.\n")
			room.Manager.RoomsCh <- room

			return
		}

		pos, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			send(player, "Invalid input: please enter a number.\n")
			continue
		}
		if pos < 1 || pos > 9 {
			send(player, "Out of range: enter a number between 1 and 9.\n")
			continue
		}

		winner, err := room.makeMove(pos, marks[idx])
		if err != nil {
			send(player, "Invalid move: %v. Try again.\n", err)
			continue
		}

		if winner != "" {
			if winner == "Draw" {
				room.BroadcastMessage("Game over! It's a draw!\n")
				room.Manager.RoomsCh <- room
				return
			}
			room.BroadcastMessage("Game over! %s wins!\n", winner)
			room.Manager.RoomsCh <- room
			return
		}

		room.PrintBoard()
		room.Turn = 1 - idx
	}
}

func (r *Room) BroadcastMessage(format string, a ...any) {
	for _, player := range r.Players {
		send(player, format, a...)
	}
}

func (r *Room) PrintBoard() {
	for _, player := range r.Players {
		engine.PrintBoard(player, r.GameBoard)
	}
}

func (r *Room) makeMove(pos int, move engine.Move) (string, error) {
	switch pos {
	case 1:
		err := r.verifyMove(0, 0)
		if err != nil {
			return "", err
		}

		r.GameBoard[0][0] = move
	case 2:
		err := r.verifyMove(0, 1)
		if err != nil {
			return "", err
		}

		r.GameBoard[0][1] = move
	case 3:
		err := r.verifyMove(0, 2)
		if err != nil {
			return "", err
		}

		r.GameBoard[0][2] = move
	case 4:
		err := r.verifyMove(1, 0)
		if err != nil {
			return "", err
		}

		r.GameBoard[1][0] = move
	case 5:
		err := r.verifyMove(1, 1)
		if err != nil {
			return "", err
		}

		r.GameBoard[1][1] = move
	case 6:
		err := r.verifyMove(1, 2)
		if err != nil {
			return "", err
		}

		r.GameBoard[1][2] = move
	case 7:
		err := r.verifyMove(2, 0)
		if err != nil {
			return "", err
		}

		r.GameBoard[2][0] = move
	case 8:
		err := r.verifyMove(2, 1)
		if err != nil {
			return "", err
		}

		r.GameBoard[2][1] = move
	case 9:
		err := r.verifyMove(2, 2)
		if err != nil {
			return "", err
		}

		r.GameBoard[2][2] = move
	}

	gameWinner := ""
	if engine.Terminal(r.GameBoard) {
		winner := engine.Winner(r.GameBoard)
		switch winner {
		case engine.EMPTY:
			gameWinner = "Draw"
		default:
			gameWinner = move.String()
		}
	}

	return gameWinner, nil
}

func (r *Room) verifyMove(row, col int) error {
	// Utility function to verify the players move
	if r.GameBoard[row][col] != engine.EMPTY {
		return errors.New("invalid move")
	}

	return nil
}
