package models

//go:generate easyjson -all

type User struct {
	ID        int
	Username  string         `json:"username"`
	Email     string         `json:"email",valid:"email"`
	Rating    int            `json:"rating"`
	AvatarUrl string         `json:"avatar_url"`
	Password  string         `json:"password,omitempty",valid:"length(8|20)"`
}
