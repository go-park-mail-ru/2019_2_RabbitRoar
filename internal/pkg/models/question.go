package models

type Question struct {
	ID     int
	Text   string
	Media  string
	Answer string
	Rating int
	Author int
	Tags   string // someday change it to slice
}
