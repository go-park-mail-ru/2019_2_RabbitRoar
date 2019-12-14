package game

type EventType string

const (
	// on websocket opening or closure
	// only for backend
	WsUpdated         EventType = "ws_run"

	// on user manually leaving the game via event
	// only for backend
	PlayerLeft        EventType = "player_left"

	// backend responding with this after ws connection status is changed
	// backend -> frontend
	UserConnected     EventType = "user_connected"

	// frontend notifies backend with this after a player is ready to start the game
	// frontend -> backend
	PlayerReadyFront  EventType = "player_ready_front"

	// backend notifies all players with this after a player readiness status is changed
	// backend -> frontend
	PlayerReadyBack   EventType = "player_ready_back"

	// backend notifies all players that the game is started after all players are ready
	// backend -> frontend
	GameStart         EventType = "start_game"

	// backend requests a question indexes from specific player and notifies all players of it
	// backend -> frontend
	RequestQuestion   EventType = "request_question_from_player"

	// frontend responses with a question indexes chosen by specific player
	// frontend -> backend
	QuestionChosen    EventType = "question_chosen"

	// backend notifies all players that respondent for the chosen question is required.
	// backend -> frontend
	RequestRespondent EventType = "request_respondent"

	// frontend notifies backend that a player is ready to answer a question.
	// frontend -> backend
	RespondentReady   EventType = "respondent_ready"

	// backend notifies all players that an answer for chosen question is required from specific player
	// backend -> frontend
	RequestAnswer     EventType = "request_answer_from_respondent"

	// frontend notifies backend with answer given by specific player
	// frontend -> backend
	AnswerGiven       EventType = "respondent_answer_given"

	// backend notifies all players of given answer with this
	// backend -> frontend
	AnswerGivenBack   EventType = "answer_given_back"

	// backend notifies host of correct answer and pends for verdict on given answer from host
	// backend -> frontend
	RequestVerdict    EventType = "request_verdict_from_host"

	// Host notifies backend that given answer is correct
	// frontend -> backend
	VerdictCorrect    EventType = "verdict_correct"

	// Host notifies backend that given answer is incorrect
	// frontend -> backend
	VerdictWrong      EventType = "verdict_wrong"

	// Backend notifies all players of given verdict and sends correct answer if the verdict if "Correct"
	// backend -> frontend
	VerdictGivenBack  EventType = "verdict_given_back"

	// Backend notifies all players that the game is ended
	// backend -> frontend
	GameEnded         EventType = "game_ended"
)

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type EventWrapper struct {
	SenderID int
	Event    *Event
}

type UserConnectedPayload struct {
	RoomName string       `json:"room_name"`
	PackName string       `json:"pack_name"`
	Host     PlayerInfo   `json:"host"`
	Players  []PlayerInfo `json:"players"`
}

type PlayerReadyBackPayload struct {
	Players []PlayerInfo `json:"players"`
}

type GameStartPayload struct {
	Themes [5]string `json:"themes"`
}

type RequestQuestionPayload struct {
	QuestionSelectorID int        `json:"player_id"`
	QuestionsStatus    [5][5]bool `json:"questions"`
}

type RequestRespondentPayload struct {
	Question    string `json:"question"`
	ThemeIdx    int    `json:"theme_id"`
	QuestionIdx int    `json:"question_id"`
}

type RequestAnswerPayload struct {
	RespondentID int `json:"player_id"`
}

type AnswerGivenBackPayload struct {
	PlayerAnswer string `json:"player_answer"`
	RespondentID int    `json:"player_id"`
}

type RequestVerdictPayload struct {
	CorrectAnswer string `json:"answer"`
}

type VerdictPayload struct {
	Verdict       bool   `json:"verdict"`
	CorrectAnswer string `json:"answer"`
}

type GameEndedPayload struct {
	Players []PlayerInfo `json:"players"`
}