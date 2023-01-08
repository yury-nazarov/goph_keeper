package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
	"github.com/yury-nazarov/goph_keeper/internal/models"

	"github.com/spf13/cobra"
)

var secretListCmd = &cobra.Command{
	Use:   "list",
	Short: "List of secret",
	Long:  `List of secret`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API
		apiServer := fmt.Sprintf("%s/api/v1/secret/list", ct.APIServer)
		httpStatus, responseBody, err := ct.HTTPClient(apiServer, http.MethodGet, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}

		// Десериализуем полученный ответ в структуру []models.Secret для дальнейшего представления
		var secrets []models.Secret
		if err = json.Unmarshal(responseBody, &secrets); err != nil {
			ct.Log.Warn(err.Error())
		}

		// Формат для пользователя в терминате
		ct.ListOfSecrets(secrets).Print()

		// Статус обработки запроса
		fmt.Println(ct.DisplayMsg(httpStatus))
	},
}

func init() {
	secretCmd.AddCommand(secretListCmd)
}


