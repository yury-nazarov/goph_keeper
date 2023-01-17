package handler

import (
	"fmt"
	"github.com/spf13/cobra"
)

var secretUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update secret by ID",
	Long:  `Update secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := App.Secret.Update(App.Item, App.Crypto)
		if err != nil {
			fmt.Printf("update secret fail: %s", err)
		}
		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	secretCmd.AddCommand(secretUpdateCmd)

	secretUpdateCmd.Flags().IntVarP(&App.Item.ID, "id", "i", 0, "id of secrets")
	secretUpdateCmd.Flags().StringVarP(&App.Item.Name, "name", "n", "", "secret_name")
	secretUpdateCmd.Flags().StringVarP(&App.Item.Data, "data", "d", "", "JSON secret data")
	secretUpdateCmd.Flags().StringVarP(&App.Item.Description, "description", "m", "", "description about secret. Optional")

	secretUpdateCmd.MarkFlagRequired("id")
}
