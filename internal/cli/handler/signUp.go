package handler

import (
	"fmt"

	"github.com/yury-nazarov/goph_keeper/internal/cli/service/auth"

	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create new account",
	Long:  `Create new account`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем все что нужно для работы метода
		a, err := auth.New()
		if err != nil {
			fmt.Printf("SignUp fail: %s", err)
		}

		// Данные полученые из флагов присваиваем в объект user и передаем дальше в слой бизнес логики
		status, err := a.SignUp(user)
		if err != nil {
			fmt.Printf("SignUp fail: %s", err)
		}

		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	Cmd.AddCommand(signUpCmd)
	signUpCmd.Flags().StringVarP(&user.Login, "login", "l", "", "username")
	signUpCmd.Flags().StringVarP(&user.Password, "password", "p", "", "pa$$w0rd")

	signUpCmd.MarkFlagRequired("login")
	signUpCmd.MarkFlagRequired("password")
}
