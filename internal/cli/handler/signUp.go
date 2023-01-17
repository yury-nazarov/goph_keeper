package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create new account",
	Long:  `Create new account`,
	Run: func(cmd *cobra.Command, args []string) {
		// Данные полученые из флагов присваиваем в объект user и передаем дальше в слой бизнес логики
		status, err := App.Auth.SignUp(App.User)
		if err != nil {
			fmt.Printf("SignUp fail: %s", err)
		}
		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	App.Cmd.AddCommand(signUpCmd)
	signUpCmd.Flags().StringVarP(&App.User.Login, "login", "l", "", "username")
	signUpCmd.Flags().StringVarP(&App.User.Password, "password", "p", "", "pa$$w0rd")

	signUpCmd.MarkFlagRequired("login")
	signUpCmd.MarkFlagRequired("password")
}
