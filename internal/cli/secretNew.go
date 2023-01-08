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

		// Данные полученые из флагов сериализуем в JSON для HTTP Request
		body, err := json.Marshal(&secret)
		if err != nil {
			fmt.Println(err)
		}

		// Запрос в HTTP API
		apiServer := fmt.Sprintf("%s/api/v1/secret/new", ct.APIServer)
		requestBody := bytes.NewBuffer(ct.Encrypt(body))
		httpStatus, _, err := ct.HTTPClient(apiServer, http.MethodPost, requestBody)
		if err != nil {
			ct.Log.Warn(err.Error())
		}

		// Статус обработки запроса
		fmt.Println(ct.DisplayMsg(httpStatus))
	},
}

func init() {
	secretCmd.AddCommand(secretNewCmd)

	secretNewCmd.Flags().StringVarP(&secret.Name, "name", "n", "", "secret_name")
	secretNewCmd.Flags().StringVarP(&secret.Data, "data", "d", "", "JSON secret data")
	secretNewCmd.Flags().StringVarP(&secret.Description, "description", "m","", "description about secret. Optional")

	secretNewCmd.MarkFlagRequired("name")
	secretNewCmd.MarkFlagRequired("data")
}


