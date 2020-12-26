package dto

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Roles    string `json:"roles"`
}

type UserId struct {
	Id int `json:"id"`
}

type ErrResp struct {
	Error string `json:"txt"`
}

type TokenRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenDTO struct {
	Token string `json:"token"`
}

type CardDTO struct {
	Number  string `json:"number"`
	Balance int    `json:"balance"`
}
