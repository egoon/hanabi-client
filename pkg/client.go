package pkg

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/egoon/hanabi-server/pkg/model"
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
	return c.connectToGame(model.Action{
		Type:         model.ActionCreate,
		ActivePlayer: model.PlayerID(c.playerName),
		GameID:       model.GameID(name),
	})
}

func (c *ClientImpl) JoinGame(name string) (Game, error) {
	return c.connectToGame(model.Action{
		Type:         model.ActionCreate,
		ActivePlayer: model.PlayerID(c.playerName),
		GameID:       model.GameID(name),
	})
}

func (c *ClientImpl) connectToGame(action model.Action) (Game, error) {
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
	state := model.GameState{}
	err = json.Unmarshal(input, &state)
	if err != nil {
		errMsg := model.Error{}
		err = json.Unmarshal(input, errMsg)
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf(errMsg.Message)
		}
	}

}
