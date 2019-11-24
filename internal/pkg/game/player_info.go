package game

type PlayerInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Ready    bool   `json:"ready"`
}
