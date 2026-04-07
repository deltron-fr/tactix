package server

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/deltron-fr/tactix/cmd/client"
	"github.com/deltron-fr/tactix/engine"
)

func startRepl(cli client.Client) {
	reader := bufio.NewReader(cli)

	fmt.Fprintln(cli, "================ Welcome to TacTix! ===============")
	fmt.Fprintln(cli, "=================================================")
	fmt.Fprintln(cli, "		1  |  2  |  3		")
	fmt.Fprintln(cli, "		-------------		")
	fmt.Fprintln(cli, "		4  |  5  |  6		")
	fmt.Fprintln(cli, "		--------------		")
	fmt.Fprintln(cli, "		7  |  8  |  9		")

	fmt.Fprintln(cli, "Game has started!")

	fmt.Fprintln(cli, "Choose a mode: ")
	fmt.Fprintln(cli, "1. Play against another player")
	fmt.Fprintln(cli, "2. Play vs AI")

	mode := ""
	for {
		fmt.Fprint(cli, "1 or 2 > ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		mode = strings.TrimSpace(line)
		if mode != "1" && mode != "2" {
			fmt.Fprintln(cli, "enter a valid choice!")
			continue
		}
		break
	}

	if mode == "2" {
		engine.PlayAI(cli)
	}
}
