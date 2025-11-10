package dto

type UserStatusResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

type LeaderboardUserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
	Rank     int    `json:"rank"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
