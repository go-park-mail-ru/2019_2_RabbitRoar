package game

type EventType string

const (
	WsRun             EventType = "ws_run"
	UserConnected     EventType = "user_connected"
	GameStart         EventType = "start_game"
	RequestQuestion   EventType = "request_question_from_player"
	QuestionChosen    EventType = "question_chosen"
	RequestRespondent EventType = "request_respondent"
	RespondentReady   EventType = "respondent_ready"
	RequestAnswer     EventType = "request_answer_from_respondent"
	AnswerGiven       EventType = "respondent_answer_given"
	AnswerGivenBack   EventType = "answer_given_back"
	RequestVerdict    EventType = "request_verdict_from_host"
	VerdictCorrect    EventType = "verdict_correct"
	VerdictWrong      EventType = "verdict_wrong"
	VerdictGivenBack  EventType = "verdict_given_back"
	GameEnded         EventType = "game_ended"
)

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type UserConnectedPayload struct {
	RoomName string       `json:"room_name"`
	PackName string       `json:"pack_name"`
	Players  []PlayerInfo `json:"players"`
}

type GameStartPayload struct {
	Themes [5]string `json:"themes"`
}

type RequestFromPlayerPayload struct {
	PlayerID int `json:"player_id"`
}

type QuestionChosenPayload struct {
	Theme       int `json:"theme"`
	QuestionIdx int `json:"question_idx"`
}

type AnswerPayload struct {
	Answer string `json:"answer"`
}

type VerictPayload struct {
	Verdict bool `json:"verdict"`
}

type EventWrapper struct {
	SenderID int
	Event    *Event
}

func NewEvent(et EventType, data ...interface{}) *Event {
	e := Event{
		Type: et,
	}

	switch et {
	case UserConnected:
		e.Payload = UserConnectedPayload{
			RoomName: data[0].(string),
			PackName: data[1].(string),
			Players:  data[2].([]PlayerInfo),
		}

	case GameStart:
		e.Payload = GameStartPayload{
			Themes: data[0].([5]string),
		}

	case RequestQuestion:
		e.Payload = RequestFromPlayerPayload{
			PlayerID: data[0].(int),
		}

	case QuestionChosen:
		e.Payload = QuestionChosenPayload{
			Theme:       data[0].(int),
			QuestionIdx: data[1].(int),
		}

	case RequestAnswer:
		e.Payload = RequestFromPlayerPayload{
			PlayerID: data[0].(int),
		}

	case AnswerGiven:
		e.Payload = AnswerPayload{
			Answer: data[0].(string),
		}

	case AnswerGivenBack:
		e.Payload = AnswerPayload{
			Answer: data[0].(string),
		}

	case RequestVerdict:
		e.Payload = RequestFromPlayerPayload{
			PlayerID: data[0].(int),
		}

	case VerdictGivenBack:
		e.Payload = VerictPayload{
			Verdict: data[0].(bool),
		}
	}

	return &e
}
