package game

import "math/rand"

type QuestionTable struct {
	questions interface{}
}

func NewQuestionTable(questions interface{}) *QuestionTable {
	return &QuestionTable{
		questions:questions,
	}
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

func (qt *QuestionTable) GetQuestion(themeID, questionID int) string {
	themeSlice := qt.questions.([]interface{})

	for themeIdx, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		if themeID == themeIdx {
			questionSlice := theme["questions"].([]interface{})
			for questionIdx, question := range questionSlice {
				question := question.(map[string]interface{})
				if questionIdx == questionID {
					return question["text"].(string)
				}
			}
		}
	}

	return ""
}

func (qt *QuestionTable) GetAnswer(themeID, questionID int) string {
	themeSlice := qt.questions.([]interface{})

	for themeIdx, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		if themeID == themeIdx {
			questionSlice := theme["questions"].([]interface{})
			for questionIdx, question := range questionSlice {
				question := question.(map[string]interface{})
				if questionIdx == questionID {
					return question["answer"].(string)
				}
			}
		}
	}

	return ""
}

// TODO: Must return available question. Err if no questions are available
func (qt *QuestionTable) GetRandAvailableQuestionIndexes() (int, int, error) {
	themeIdx := rand.Int() % 5
	questionIdx := rand.Int() % 5

	return themeIdx, questionIdx, nil
}