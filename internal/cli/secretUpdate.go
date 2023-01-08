package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
	"github.com/yury-nazarov/goph_keeper/internal/models"
	"io/ioutil"
	"net/http"
)

var secretUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update secret by ID",
	Long:  `Update secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API
		// 1. Получаем секрет по ID
		resp, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/secret/%d", ct.APIServer, secret.ID), http.MethodGet, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// JSON из набора байт
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		var originSecret models.Secret
		if err = json.Unmarshal(body, &originSecret); err != nil {
			ct.Log.Warn(err.Error())
		}

		// Формат для пользователя в терминате
		fmt.Printf("Origin  secret: %+v\n", originSecret)


		// Для секрета полученого из БД обновляем поля: Name, Data, Description - если они были изменены и не пустые.
		if originSecret.Name != secret.Name && len(secret.Name) > 0 {
			originSecret.Name = secret.Name
		}
		if originSecret.Data != secret.Data && len(secret.Data) > 0  {
			originSecret.Data = secret.Data
		}
		if originSecret.Description != secret.Description && len(secret.Description) > 0 {
			originSecret.Description = secret.Description
		}
		// Результат
		fmt.Printf("Updated secret: %+v\n", originSecret)

		// Сериализуем в JSON отправлеяем в HTTP API
		// JSON для HTTP Request
		body, err = json.Marshal(&originSecret)
		if err != nil {
			fmt.Println(err)
		}

		// Запрос в HTTP API
		resp, err = ct.HTTPClient(fmt.Sprintf("%s/api/v1/secret/update", ct.APIServer), http.MethodPut, bytes.NewBuffer(ct.Encrypt(body)))
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// Вывод в терминал статус операции
		fmt.Println(ct.DisplayMsg(resp.Status))
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


