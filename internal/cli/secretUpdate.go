package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
	"github.com/yury-nazarov/goph_keeper/internal/models"

	"github.com/spf13/cobra"
)

var secretUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update secret by ID",
	Long:  `Update secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API. Получаем секрет по ID
		apiServer := fmt.Sprintf("%s/api/v1/secret/%d", ct.APIServer, secret.ID)
		_, responseBody, err := ct.HTTPClient(apiServer, http.MethodGet, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}

		// Десериализуем полученный ответ в структуру models.Secret для дальнейшего представления
		var originSecret models.Secret
		if err = json.Unmarshal(responseBody, &originSecret); err != nil {
			ct.Log.Warn(err.Error())
		}

		// secret - секрет измененный пользователей
		// originSecret - секрет полученый из БД.
		// Для него обновляем поля:
		//			Name, Data, Description - если они были изменены и не пустые.
		if originSecret.Name != secret.Name && len(secret.Name) > 0 {
			originSecret.Name = secret.Name
		}
		if originSecret.Data != secret.Data && len(secret.Data) > 0  {
			originSecret.Data = secret.Data
		}
		if originSecret.Description != secret.Description && len(secret.Description) > 0 {
			originSecret.Description = secret.Description
		}

		// Шифруем секрет
		originSecret.Data = ct.Encrypt([]byte(originSecret.Data))
		// Сериализуем originSecret в JSON отправлеяем в HTTP API
		body, err := json.Marshal(&originSecret)
		if err != nil {
			fmt.Println(err)
		}
		// Запрос в HTTP API для обновления данных о секрете
		apiServer = fmt.Sprintf("%s/api/v1/secret/update", ct.APIServer)
		httpStatus, _, err  := ct.HTTPClient(apiServer, http.MethodPut, bytes.NewBuffer(body))
		if err != nil {
			ct.Log.Warn(err.Error())
		}

		// Статус обработки запроса
		fmt.Println(ct.DisplayMsg(httpStatus))
	},
}

func init() {
	secretCmd.AddCommand(secretUpdateCmd)

	secretUpdateCmd.Flags().IntVarP(&secret.ID, "id", "i",0, "id of secrets")
	secretUpdateCmd.Flags().StringVarP(&secret.Name, "name", "n", "", "secret_name")
	secretUpdateCmd.Flags().StringVarP(&secret.Data, "data", "d", "", "JSON secret data")
	secretUpdateCmd.Flags().StringVarP(&secret.Description, "description", "m","", "description about secret. Optional")

	secretUpdateCmd.MarkFlagRequired("id")
}


