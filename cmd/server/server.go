package server

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

	"github.com/deltron-fr/tactix/cmd/client"
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
	GameBoard Board
	Players   [2]client.Client
	Turn      int
	Mu        sync.Mutex
}

func generateRoomID() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return "ROOM-" + string(b)
}

func NewRoom(player1, player2 client.Client) *Room {
	return &Room{
		ID:        generateRoomID(),
		GameBoard: make(Board, 3),
		Players:   [2]client.Client{player1, player2},
		Turn:      0,
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

		handleMatchMaking(cli)
	}
}

func handleMatchMaking(cli client.Client) {
	fmt.Printf("New client connected: %v\n", cli.RemoteAddr())
	cli.Write([]byte("Welcome to TacTix! Waiting for an opponent...\n"))
	matchMaker <- cli

	if len(matchMaker) < 2 {
		return
	}

	var player1, player2 client.Client

	player1 = <-matchMaker
	player2 = <-matchMaker

	room := NewRoom(player1, player2)
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
	marks := [2]Move{X, O}

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

		if err := room.makeMove(pos, marks[idx]); err != nil {
			send(player, "Invalid move: %v. Try again.\n", err)
			continue
		}

		room.Turn = 1 - idx
	}
}

func (r *Room) makeMove(pos int, move Move) error {
	switch pos {
	case 1:
		err := r.verifyMove(0, 0)
		if err != nil {
			return err
		}

		r.GameBoard[0][0] = move
	case 2:
		err := r.verifyMove(0, 1)
		if err != nil {
			return err
		}

		r.GameBoard[0][1] = move
	case 3:
		err := r.verifyMove(0, 2)
		if err != nil {
			return err
		}

		r.GameBoard[0][2] = move
	case 4:
		err := r.verifyMove(1, 0)
		if err != nil {
			return err
		}

		r.GameBoard[1][0] = move
	case 5:
		err := r.verifyMove(1, 1)
		if err != nil {
			return err
		}

		r.GameBoard[1][1] = move
	case 6:
		err := r.verifyMove(1, 2)
		if err != nil {
			return err
		}

		r.GameBoard[1][2] = move
	case 7:
		err := r.verifyMove(2, 0)
		if err != nil {
			return err
		}

		r.GameBoard[2][0] = move
	case 8:
		err := r.verifyMove(2, 1)
		if err != nil {
			return err
		}

		r.GameBoard[2][1] = move
	case 9:
		err := r.verifyMove(2, 2)
		if err != nil {
			return err
		}

		r.GameBoard[2][2] = move
	}

	return nil
}

func (r *Room) verifyMove(row, col int) error {
	// Utility function to verify the players move
	if r.GameBoard[row][col] != EMPTY {
		return errors.New("invalid move")
	}

	return nil
}
