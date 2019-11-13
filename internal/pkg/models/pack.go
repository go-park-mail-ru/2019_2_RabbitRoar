package models

type Pack struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description, omitempty"`
	Rating      int         `json:"rating"`
	Author      int         `json:"author"`
	Tags        string      `json:"tags"`
	Questions   interface{} `json:"pack,omitempty"`
}
