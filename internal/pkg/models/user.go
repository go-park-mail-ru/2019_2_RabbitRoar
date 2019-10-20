package models

type User struct {
	UID       int64
	Username  string `json:"username"`
	Email     string `json:"email",valid:"email"`
	Rating    int64  `json:"rating"`
	AvatarUrl string `json:"avatar_url"`
	Password  string `json:",omitempty",valid:"length(8|20)"`
}
