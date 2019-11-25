package pack

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Sanitizer interface {
	Sanitize(p *models.Pack) models.Pack
	SanitizeSlice(pSlice []models.Pack) []models.Pack
}
