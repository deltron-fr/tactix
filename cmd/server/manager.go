package main

import "time"

type Manager struct {
	RoomsCh chan *Room
}

func NewManager() *Manager {
	return &Manager{
		RoomsCh: make(chan *Room),
	}
}

func (m *Manager) handleRoomCleanup() {
	for room := range m.RoomsCh {
		for _, player := range room.Players {
			send(player, "The Game has come to an end. Thanks for playing!\n")
			time.Sleep(100 * time.Millisecond)
			player.Close()
		}
	}
}
