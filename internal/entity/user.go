package entity

type User struct {
	UID      int64 `json:"id"`
	Name     string
	Password string
	Rating   int64
}
