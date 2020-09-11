package model

import "github.com/egoon/hanabi-server/pkg/model"

type GameMsg struct {
	State *model.GameState
	Err *model.Error
}
