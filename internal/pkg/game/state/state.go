package state

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"

type pendPlayers struct {
}

func (state *pendPlayers) handle(e game.Event) *game.State {
	return nil
}

type pendQuestionChosen struct {
}

func (state *pendQuestionChosen) handle(e game.Event) *game.State {
	return nil
}

type pendRespondent struct {
}

func (state *pendRespondent) handle(e game.Event) *game.State {
	return nil
}

type pendAnswer struct {
}

func (state *pendAnswer) handle(e game.Event) *game.State {
	return nil
}

type pendVerdict struct {
}

func (state *pendVerdict) handle(e game.Event) *game.State {
	return nil
}
