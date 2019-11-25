package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/microcosm-cc/bluemonday"
)

type packSanitizer struct {
	sanitizer *bluemonday.Policy
}

func NewPackSanitizer(policy *bluemonday.Policy) pack.Sanitizer {
	return &packSanitizer{
		sanitizer: policy,
	}
}

//TODO: make me more safe (or contract that DB has valid form)
func (s *packSanitizer) sanitizeQuestions(p interface{}) {
	if p == nil {
		return
	}

	themeSlice := p.([]interface{})

	for _, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		theme["name"] = s.sanitizer.Sanitize(theme["name"].(string))
		questionSlice := theme["questions"].([]interface{})
		for _, question := range questionSlice {
			question := question.(map[string]interface{})
			question["text"] = s.sanitizer.Sanitize(question["text"].(string))
			question["answer"] = s.sanitizer.Sanitize(question["answer"].(string))
		}
	}
}

func (s *packSanitizer) Sanitize(p *models.Pack) models.Pack {
	p.Name = s.sanitizer.Sanitize(p.Name)
	p.Description = s.sanitizer.Sanitize(p.Description)
	p.Tags = s.sanitizer.Sanitize(p.Tags)
	s.sanitizeQuestions(p.Questions)
	return *p
}

func (s *packSanitizer) SanitizeSlice(p []models.Pack) []models.Pack {
	for i := 0; i < len(p); i++ {
		p[i] = s.Sanitize(&p[i])
	}
	return p
}
