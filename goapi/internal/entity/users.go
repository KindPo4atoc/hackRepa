package entity

type UserData struct {
	IdUser       string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"pass_hash"`
	Role         int    `json:"role"`
	Email        string `json:"email"`
}
