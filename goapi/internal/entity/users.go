package entity

type UserData struct {
	IdUser       string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
	Role         int    `json:"role"`
}
