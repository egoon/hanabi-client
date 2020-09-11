package pkg

import (
	"encoding/json"
	"fmt"
	server "github.com/egoon/hanabi-server/pkg/model"
	"net"
)

const pingInterval = 10 //secods

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

	actionJson, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(actionJson)
	if err != nil {
		conn.Close()
		return nil, err
	}
	var input []byte
	conn.Read(input)
	state := server.GameState{}
	err = json.Unmarshal(input, &state)
	if err != nil {
		errMsg := server.Error{}
		err = json.Unmarshal(input, errMsg)
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf(errMsg.Message)
		}
	}
	return NewGame(conn), nil
}
