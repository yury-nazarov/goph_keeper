package handler

import (
	"fmt"

	"github.com/yury-nazarov/goph_keeper/internal/cli/service/auth"

	"github.com/spf13/cobra"
)

var signInCmd = &cobra.Command{
	Use:   "signin",
	Short: "LogIn",
	Long:  `LogIn`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем все что нужно для работы метода
		a, err := auth.New()
		if err != nil {
			fmt.Printf("SignUp fail: %s", err)
		}
		status, err := a.SignIn(user)
		if err != nil {
			fmt.Printf("SignIn fail: %s", err)
		}

		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {

	Cmd.AddCommand(signInCmd)
	signInCmd.Flags().StringVarP(&user.Login, "login", "l", "", "username")
	signInCmd.Flags().StringVarP(&user.Password, "password", "p", "", "pa$$w0rd")

	signInCmd.MarkFlagRequired("login")
	signInCmd.MarkFlagRequired("password")
}
