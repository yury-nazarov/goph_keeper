package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signInCmd = &cobra.Command{
	Use:   "signin",
	Short: "LogIn",
	Long:  `LogIn`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := App.Auth.SignIn(App.User)
		if err != nil {
			fmt.Printf("SignIn fail: %s", err)
		}
		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {

	App.Cmd.AddCommand(signInCmd)
	signInCmd.Flags().StringVarP(&App.User.Login, "login", "l", "", "username")
	signInCmd.Flags().StringVarP(&App.User.Password, "password", "p", "", "pa$$w0rd")

	signInCmd.MarkFlagRequired("login")
	signInCmd.MarkFlagRequired("password")
}
