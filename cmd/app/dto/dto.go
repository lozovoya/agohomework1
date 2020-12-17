package dto

import "github.com/google/uuid"

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

type Token struct {
	Token uuid.UUID `json:"token"`
}
