package main

import "os"

const (
	port = ":8000"
)

func main() {
	initClient := os.Stdout
	startRepl(initClient)
}
