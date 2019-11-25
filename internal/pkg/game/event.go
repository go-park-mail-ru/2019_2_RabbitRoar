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
	PlayerReadyFront  EventType = "player_ready_front"
	PlayerReadyBack   EventType = "player_ready_back"
)

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type PlayerReadyBackPayload struct {
	Players []PlayerInfo `json:"players"`
}

type RequestQuestionPayload struct {
	PlayerID int `json:"player_id"`
}

type RequestAnswerPayload struct {
	PlayerID int `json:"player_id"`
}

type AnswerGivenPayload struct {
	Answer string `json:"answer"`
}

type AnswerGivenBackPayload struct {
	PlayerAnswer string `json:"player_answer"`
	PlayerID     int    `json:"player_id"`
}

type RequestVerdictPayload struct {
	HostID int    `json:"host_id"`
	Answer string `json:"answer"`
}

type RequestRespondentPayload struct {
	Question   string `json:"question"`
	ThemeID    int    `json:"theme_id"`
	QuestionID int    `json:"question_id"`
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
	ThemeIdx    int `json:"theme_idx"`
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
