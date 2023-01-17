package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signOutCmd = &cobra.Command{
	Use:   "signout",
	Short: "Logout user",
	Long:  `Logout user`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := App.Auth.SignOut()
		if err != nil {
			fmt.Printf("SignOut fail: %s", err)
		}
		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	App.Cmd.AddCommand(signOutCmd)
}
