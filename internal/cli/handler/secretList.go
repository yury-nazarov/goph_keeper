package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var secretListCmd = &cobra.Command{
	Use:   "list",
	Short: "List of secret",
	Long:  `List of secret`,
	Run: func(cmd *cobra.Command, args []string) {
		status, secrets, err := App.Secret.List()
		if err != nil {
			fmt.Printf("create secret fail: %s", err)
		}

		// Выводим секреты
		App.Secret.ListOfSecrets(secrets).Print()

		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	secretCmd.AddCommand(secretListCmd)
}
