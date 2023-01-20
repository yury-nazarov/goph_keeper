package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var secretGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get of secret by ID",
	Long:  `Get of secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		status, secret, err := App.Secret.Get(App.Item.ID, App.Crypto)
		if err != nil {
			fmt.Printf("get secret fail: %s", err)
		}

		// Выводим секреты
		App.Secret.ListOfSecrets(secret).Print()

		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	secretCmd.AddCommand(secretGetCmd)
	secretGetCmd.Flags().IntVarP(&App.Item.ID, "id", "i", 0, "id of secrets")

	secretGetCmd.MarkFlagRequired("id")
}
