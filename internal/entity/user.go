package entity

type User struct {
	UID      int64
	Username string `json:"username"`
	Email    string `json:"email"`
	Rating   int64  `json:"rating"`
	Url      string `json:"url"`
	Password string
}
