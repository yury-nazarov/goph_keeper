package handler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var secretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete of secret by ID",
	Long:  `Delete of secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := App.Secret.Delete(App.Item.ID)
		if err != nil {
			fmt.Printf("delete secret fail: %s", err)
		}
		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	secretCmd.AddCommand(secretDeleteCmd)
	secretDeleteCmd.Flags().IntVarP(&App.Item.ID, "id", "i", 0, "id of secrets")

	secretDeleteCmd.MarkFlagRequired("id")
}
