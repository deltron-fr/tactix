package main

import (
	"github.com/deltron-fr/tactix/internal/engine"
)


type Game struct {
	Clients map[*Client]bool
	GameState *engine.Config
}
