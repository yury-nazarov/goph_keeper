package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"

	"github.com/spf13/cobra"
)

var secretNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new secret",
	Long:  `Create new secret`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// JSON для HTTP Request
		body, err := json.Marshal(&secret)
		if err != nil {
			fmt.Println(err)
		}

		// Запрос в HTTP API
		resp, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/secret/new", ct.APIServer), http.MethodPost, bytes.NewBuffer(ct.Encrypt(body)))
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// Вывод в терминал
		fmt.Println(ct.AuthDisplayMsg(resp.Status))
	},
}

func init() {
	secretCmd.AddCommand(secretNewCmd)

	secretNewCmd.Flags().StringVarP(&secret.Name, "name", "n", "", "secret_name")
	secretNewCmd.Flags().StringVarP(&secret.Data, "data", "d", "", "JSON secret data")
	secretNewCmd.Flags().StringVarP(&secret.Description, "message", "m","", "description about secret. Optional")

	secretNewCmd.MarkFlagRequired("name")
	secretNewCmd.MarkFlagRequired("data")
}


