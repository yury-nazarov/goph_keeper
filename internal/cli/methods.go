package cli

import (
	"fmt"
)

// AuthDisplayMsg выводит в терминал сообщение пользователю
func AuthDisplayMsg(httpStatus string) {
	switch httpStatus {
	case "200 OK":
		fmt.Printf("You are login as: %s\nPlease type flag: -h for more options.\n", user.Login)
	case "201 Created":
		fmt.Printf("Success registration. You are login as: %s\nPlease type flag: -h for more options.\n", user.Login)
	case "400 Bad Request":
		fmt.Println("Request format error")
	case "401 Unauthorized":
		fmt.Println("Incorrect login or password")
	case "409 Conflict":
		fmt.Println("Login is already exist")
	case "500 Internal Server Error":
		fmt.Println("Internal Server Error")
	default:
		fmt.Println("Something wrong. Please try later")
	}
}


func debug(status string, token string) {
	fmt.Println("DEBUG. Response status:", status)
	fmt.Println("DEBUG. Token:", token)
}

//func SaveAuth(user models.User) {
//
//	// Выглядит как решение в лоб. Но пока я не смог убрать
//	// эти данные в контекст или какую либо глобальную структурку.
//	if err = os.Setenv("GK_CLI_TOKEN", user.Token); err != nil {
//		fmt.Println("can't save user.Token")
//		os.Exit(1)
//	}
//	if err = os.Setenv("GK_CLI_USER", user.Login); err != nil {
//		fmt.Println("can't user.Login ")
//		os.Exit(1)
//	}
//}