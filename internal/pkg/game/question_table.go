package game

import (
	"errors"
)

type QuestionTable struct {
	questionsAvailable [5][5]bool
	questions          interface{}
}

func NewQuestionTable(questions interface{}) *QuestionTable {
	return &QuestionTable{
		questionsAvailable: [5][5]bool{
			{true, true, true, true, true},
			{true, true, true, true, true},
			{true, true, true, true, true},
			{true, true, true, true, true},
			{true, true, true, true, true},
		},
		questions:questions,
	}
}

func (qt *QuestionTable) IsAnyQuestionAvailable() bool {
	for _, qArr := range qt.questionsAvailable {
		for _, q := range qArr {
			if q {
				return true
			}
		}
	}
	return false
}

func (qt *QuestionTable) IsQuestionAvailable(themeIdx, questionIdx int) bool {
	return qt.questionsAvailable[themeIdx][questionIdx]
}

func (qt *QuestionTable) SetQuestionUnavailable(themeIdx, questionIdx int) {
	qt.questionsAvailable[themeIdx][questionIdx] = false
}

func (qt *QuestionTable) GetThemes() [5]string {
	var themes [5]string
	themeSlice := qt.questions.([]interface{})

	for i := 0; i < 5; i++ {
		theme := themeSlice[i].(map[string]interface{})
		themes[i] = theme["name"].(string)
	}

	return themes
}

func (qt *QuestionTable) GetQuestion(themeIdx, questionIdx int) string {
	themeSlice := qt.questions.([]interface{})

	for themeIter, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		if themeIter == themeIdx {
			questionSlice := theme["questions"].([]interface{})
			for questionIter, question := range questionSlice {
				question := question.(map[string]interface{})
				if questionIter == questionIdx {
					return question["text"].(string)
				}
			}
		}
	}

	return ""
}

func (qt *QuestionTable) GetAnswer(themeIdx, questionIdx int) string {
	themeSlice := qt.questions.([]interface{})

	for themeIter, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		if themeIter == themeIdx {
			questionSlice := theme["questions"].([]interface{})
			for questionIter, question := range questionSlice {
				question := question.(map[string]interface{})
				if questionIter == questionIdx {
					return question["answer"].(string)
				}
			}
		}
	}

	return ""
}

func (qt *QuestionTable) GetAnyAvailableQuestionIndexes() (int, int, error) {
	for themeIdx, qArr := range qt.questionsAvailable {
		for questionIdx, q := range qArr {
			if q {
				return themeIdx, questionIdx, nil
			}
		}
	}

	return 0, 0, errors.New("no question is available")
}