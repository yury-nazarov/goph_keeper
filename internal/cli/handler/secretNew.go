package handler

import (
	"fmt"
	"github.com/spf13/cobra"
)

var secretNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new secret",
	Long:  `Create new secret`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := App.Secret.New(App.Item)
		if err != nil {
			fmt.Printf("create secret fail: %s", err)
		}
		// Статус обработки запроса
		fmt.Println(status)
	},
}

func init() {
	secretCmd.AddCommand(secretNewCmd)

	secretNewCmd.Flags().StringVarP(&App.Item.Name, "name", "n", "", "secret_name")
	secretNewCmd.Flags().StringVarP(&App.Item.Data, "data", "d", "", "JSON secret data")
	secretNewCmd.Flags().StringVarP(&App.Item.Description, "description", "m", "", "description about secret. Optional")

	secretNewCmd.MarkFlagRequired("name")
	secretNewCmd.MarkFlagRequired("data")
}
