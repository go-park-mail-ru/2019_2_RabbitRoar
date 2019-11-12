package event

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"

func NewEvent(et game.EventType, data ...interface{}) *game.Event {
	e := game.Event{
		Type: et,
	}

	switch et {
	case game.GameStart:
		e.Payload = game.GameStartPayload{
			Themes: data[0].([5]string),
		}

	case game.RequestQuestion:
		e.Payload = game.RequestFromPlayerPayload{
			PlayerID: data[0].(int),
		}

	case game.QuestionChosen:
		e.Payload = game.QuestionChosenPayload{
			Theme:       data[0].(int),
			QuestionIdx: data[1].(int),
		}

	case game.RequestAnswer:
		e.Payload = game.RequestFromPlayerPayload{
			PlayerID: data[0].(int),
		}

	case game.AnswerGiven:
		e.Payload = game.AnswerPayload{
			Answer: data[0].(string),
		}

	case game.AnswerGivenBack:
		e.Payload = game.AnswerPayload{
			Answer: data[0].(string),
		}

	case game.RequestVerdict:
		e.Payload = game.RequestFromPlayerPayload{
			PlayerID: data[0].(int),
		}

	case game.VerictGivenBack:
		e.Payload = game.VerictPayload{
			Verdict: data[0].(bool),
		}
	}

	return &e
}
