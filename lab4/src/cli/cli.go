package cli

import (
	"flag"
	"log"
)

type ProgArgs struct {
	IsRobot, CreateGame       bool
	GSocket, SSocket, CSocket string
}

func Parse() *ProgArgs {
	robot := flag.Bool("robot", false, "Enable robot mode (default: false)")
	createGame := flag.Bool("create-game", false, "Create a game (default: false)")
	gsocket := flag.String("gsocket", "localhost:8080", "Socket for the game (default: localhost:8080)")
	ssocket := flag.String("ssocket", "", "Socket for the internal desktop server")
	csocket := flag.String("csocket", "", "Socket for the internal desktop client")
	flag.Parse()
	if !*robot && *createGame {
		log.Fatal("Error: only robots can enable automatic game creation.")
	}
	if *robot && (*ssocket != "" || *csocket != "") {
		log.Fatal("Error: --ssocket and --csocket cannot be used when --robot is true.")
	}

	return &ProgArgs{*robot, *createGame, *gsocket, *ssocket, *csocket}
}
