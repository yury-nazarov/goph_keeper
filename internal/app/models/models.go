package models

// User структура данных описывающая все необходимое для работы с пользователем
type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"-"`
}
