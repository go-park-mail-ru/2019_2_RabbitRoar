package game

type GameEndedState struct {
	BaseState
}

func NewGameEndedState(g *Game, ctx *StateContext) State {
	e := Event{
		Type:    GameEnded,
		Payload: GameEndedPayload{
			Players: g.GatherPlayersInfo(),
		},
	}

	g.BroadcastEvent(e)
	
	return &GameEndedState{
		BaseState: BaseState{
			Game: g,
			Ctx:  ctx,
		},
	}
}

func (s *GameEndedState) Handle(ew EventWrapper) State {
	s.Game.logger.Info("GameEnded: got event: ", ew)
	s.Game.logger.Info("GameEnded: stopping game.")
	return nil
}