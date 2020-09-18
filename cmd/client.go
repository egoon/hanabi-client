package main

import (
	"github.com/egoon/hanabi-client/pkg/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	client := model.NewClient("localhost", "testPlayer")
	game, err := client.CreateGame("testGame")
	if err != nil {
		log.Error("failed to create game", err)
	}
	messages := game.GetMessages()
	msg := <- messages
	log.Info(msg)
}

