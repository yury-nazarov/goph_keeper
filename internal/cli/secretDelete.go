package cli

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
)

var secretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete of secret by ID",
	Long:  `Delete of secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API
		resp, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/secret/delete/%d", ct.APIServer, id), http.MethodDelete, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// Вывод в терминал
		fmt.Println(ct.DisplayMsg(resp.Status))
	},
}

func init() {
	secretCmd.AddCommand(secretDeleteCmd)
	secretDeleteCmd.Flags().IntVarP(&id, "id", "i",0, "id of secrets")

	secretDeleteCmd.MarkFlagRequired("id")
}


