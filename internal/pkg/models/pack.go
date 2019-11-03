package models

type Pack struct {
	ID          int
	Name        string
	Description string
	Img         string
	Rating      int
	Author      int
	Private     bool
	Tags        string     // someday change it to slice
}
