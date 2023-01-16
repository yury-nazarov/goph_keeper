package handler

//import (
//	"encoding/json"
//	"fmt"
//	"net/http"
//
//	"github.com/yury-nazarov/goph_keeper/internal/models"
//
//	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
//
//	"github.com/spf13/cobra"
//)
//
//var secretGetCmd = &cobra.Command{
//	Use:   "get",
//	Short: "Get of secret by ID",
//	Long:  `Get of secret by ID`,
//	Run: func(cmd *cobra.Command, args []string) {
//		// Инициируем вспомогательную структуру
//		ct := tools.New()
//
//		// Запрос в HTTP API
//		apiServer := fmt.Sprintf("%s/api/v1/secret/%d", ct.APIServer, id)
//		httpStatus, responseBody, err := ct.HTTPClient(apiServer, http.MethodGet, nil)
//		if err != nil {
//			ct.Log.Warn(err.Error())
//		}
//
//		// Десериализуем полученный ответ в структуру models.Secret для дальнейшего представления
//		if err = json.Unmarshal(responseBody, &secret); err != nil {
//			ct.Log.Warn(err.Error())
//		}
//
//		// Формат для пользователя в терминате
//		secrets := []models.Secret{secret}
//		ct.ListOfSecrets(secrets).Print()
//
//		// Статус обработки запроса
//		fmt.Println(ct.DisplayMsg(httpStatus))
//	},
//}
//
//func init() {
//	secretCmd.AddCommand(secretGetCmd)
//	secretGetCmd.Flags().IntVarP(&id, "id", "i", 0, "id of secrets")
//
//	secretGetCmd.MarkFlagRequired("id")
//}
