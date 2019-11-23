package models

//go:generate easyjson -all

type Session struct {
	ID string
	User User
}
