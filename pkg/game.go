package pkg

import (
	"encoding/json"
	"github.com/egoon/hanabi-client/pkg/model"
	server "github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

const pingInterval = 10 * time.Second

type Game interface {
	GetMessages() chan model.GameMsg
}

type GameImpl struct {
	messages chan model.GameMsg
	actions chan server.Action
	conn net.Conn
	connected bool
}

func NewGame(conn net.Conn) Game {
	game := &GameImpl{
		messages: make(chan model.GameMsg, 5),
		actions: make(chan server.Action, 5),
		conn: conn,
		connected: true,
	}
	go game.sendActions()
	go game.readResponse()

	return game
}

func (g *GameImpl) GetMessages() chan model.GameMsg {
	return g.messages
}

func (g *GameImpl) sendActions() {
	defer g.conn.Close()
	ping, err := json.Marshal(server.Action{
		Type:         server.ActionPing,
	})
	if err != nil {
		log.Fatal("failed to marshal ping action", err)
	}
	for g.connected {
		select {
		case action := <- g.actions:
			actionJson, err := json.Marshal(action)
			if err != nil {
				log.Error("failed to marshal action", action, err)
			} else {
				g.send(actionJson)
			}
		// routine will ping server every 10 seconds
		case <- time.After(pingInterval):
			log.Info("sending ping")
			g.send(ping)
		}
	}
}

func (g *GameImpl) send(ping []byte) {
	_, err := g.conn.Write(ping)
	if err != nil {
		log.Error("failed to send action to server", err)
		g.connected = false
	}
}

func (g *GameImpl) readResponse() {
	defer g.conn.Close()
	readBuffer := make([]byte, 1000)
	for g.connected {
		bytesRead, err := g.conn.Read(readBuffer)
		if err != nil {
			log.Error("read from server failed", err)
		}
		state := server.GameState{}
		err = json.Unmarshal(readBuffer[:bytesRead], &state)
		if err == nil {
			g.messages <- model.GameMsg{State: &state}
		} else {
			msg := server.Error{}
			err = json.Unmarshal(readBuffer[:bytesRead], &msg)
			if err == nil {
				g.messages <-  model.GameMsg{Err: &msg}
			} else {
				log.Error("failed to parse message from server: ", string(readBuffer[:bytesRead]))
			}
		}
	}
}