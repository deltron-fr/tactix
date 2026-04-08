package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/deltron-fr/tactix/engine"
)

func startRepl(cli io.ReadWriter) {
	reader := bufio.NewReader(cli)

	fmt.Fprintln(cli, "================ Welcome to TacTix! ===============")
	fmt.Fprintln(cli, "=================================================")
	fmt.Fprintln(cli, "		1  |  2  |  3		")
	fmt.Fprintln(cli, "		-------------		")
	fmt.Fprintln(cli, "		4  |  5  |  6		")
	fmt.Fprintln(cli, "		--------------		")
	fmt.Fprintln(cli, "		7  |  8  |  9		")

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
	} else {
		conn, err := net.Dial("tcp", "localhost"+port)
		if err != nil {
			log.Fatal("ERROR - Error connecting to server:", err)
		}
		defer conn.Close()

		go func() {
			mustCopy(cli, conn)

			fmt.Fprintln(cli, "Connection closed by server. Exiting...")
			os.Exit(0)
		}()
		mustCopy(conn, os.Stdin)
	}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal("ERROR - Error copying data:", err)
	}
}
