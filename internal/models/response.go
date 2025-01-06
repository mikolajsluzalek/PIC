package models

type LoginResponse struct {
	JWT string `json:"jwt"`
	Exp int64  `json:"exp"`
}
