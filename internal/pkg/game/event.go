package game

type EventType string

const (
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
	VerictGivenBack   EventType = "verdict_given_back"
)

type Event struct {
	Type    EventType
	Payload interface{}
}

type GameStartPayload struct {
	Themes [5]string
}

type RequestFromPlayerPayload struct {
	PlayerID int
}

type QuestionChosenPayload struct {
	Theme       int
	QuestionIdx int
}

type AnswerPayload struct {
	Answer string
}

type VerictPayload struct {
	Verdict bool
}
