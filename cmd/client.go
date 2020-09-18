package main

import (
	"github.com/egoon/hanabi-client/pkg"
	log "github.com/sirupsen/logrus"
)

func main() {
	client := pkg.NewClient("localhost", "testPlayer")
	game, err := client.CreateGame("testGame")
	if err != nil {
		log.Error("failed to create game", err)
	}
	messages := game.GetMessages()
	msg := <- messages
	log.Info(msg)
}

