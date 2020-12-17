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
