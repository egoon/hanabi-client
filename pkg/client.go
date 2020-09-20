package pkg

import (
	"encoding/json"
	"fmt"
	server "github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
	"net"
)

func NewClient(url string, playerName string) Client {
	return &ClientImpl{
		url:        url,
		playerName: playerName,
	}
}

type Client interface {
	CreateGame(name string) (Game, error)
	JoinGame(name string) (Game, error)
}

type ClientImpl struct {
	url        string
	playerName string
}

func (c *ClientImpl) CreateGame(name string) (Game, error) {
	return c.connectToGame(server.Action{
		Type:         server.ActionCreate,
		ActivePlayer: server.PlayerID(c.playerName),
		GameID:       server.GameID(name),
	})
}

func (c *ClientImpl) JoinGame(name string) (Game, error) {
	return c.connectToGame(server.Action{
		Type:         server.ActionCreate,
		ActivePlayer: server.PlayerID(c.playerName),
		GameID:       server.GameID(name),
	})
}

func (c *ClientImpl) connectToGame(action server.Action) (Game, error) {
	conn, err := net.Dial("tcp", c.url+":579")
	if err != nil {
		return nil, err
	}
	log.Info("Connected")

	actionJson, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	log.Info("Action: ", string(actionJson))
	_, err = conn.Write(actionJson)
	if err != nil {
		conn.Close()
		return nil, err
	}
	input := make([]byte, 1000)
	bytesRead, err := conn.Read(input)
	log.Info("state: ", string(input[:bytesRead]))
	if err != nil {
		return nil, err
	}
	state := server.GameState{}
	err = json.Unmarshal(input[:bytesRead], &state)
	if err != nil {
		log.Debug("parsing state failed: ", err)
		errMsg := server.Error{}
		err = json.Unmarshal(input[:bytesRead], &errMsg)
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf(errMsg.Message)
		}
	}
	return NewGame(conn), nil
}
