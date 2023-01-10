package models

// User структура данных описывающая все необходимое для работы с пользователем
type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"-"`
}

// Secret описывает структуру данных для работы с секретами
type Secret struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Data        string `json:"data"`
	Description string `json:"description"`
}
