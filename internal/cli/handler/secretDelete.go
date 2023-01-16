package handler

//import (
//	"fmt"
//	"net/http"
//
//	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
//
//	"github.com/spf13/cobra"
//)
//
//var secretDeleteCmd = &cobra.Command{
//	Use:   "delete",
//	Short: "Delete of secret by ID",
//	Long:  `Delete of secret by ID`,
//	Run: func(cmd *cobra.Command, args []string) {
//		// Инициируем вспомогательную структуру
//		ct := tools.New()
//
//		// Запрос в HTTP API
//		apiServer := fmt.Sprintf("%s/api/v1/secret/delete/%d", ct.APIServer, id)
//		httpStatus, _, err := ct.HTTPClient(apiServer, http.MethodDelete, nil)
//		if err != nil {
//			ct.Log.Warn(err.Error())
//		}
//
//		// Статус обработки запроса
//		fmt.Println(ct.DisplayMsg(httpStatus))
//	},
//}
//
//func init() {
//	secretCmd.AddCommand(secretDeleteCmd)
//	secretDeleteCmd.Flags().IntVarP(&id, "id", "i", 0, "id of secrets")
//
//	secretDeleteCmd.MarkFlagRequired("id")
//}
