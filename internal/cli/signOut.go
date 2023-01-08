package cli

import (
	"fmt"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"

	"github.com/spf13/cobra"
)

var signOutCmd = &cobra.Command{
	Use:   "signout",
	Short: "Logout user",
	Long:  `Logout user`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API для удаления сессии
		httpStatus, _, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/auth/signout", ct.APIServer), http.MethodDelete, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}

		// Статус обработки запроса
		fmt.Println(ct.DisplayMsg(httpStatus))

		// Удаляем временный файл
		ct.AuthDel()
	},
}

func init() {
	Cmd.AddCommand(signOutCmd)
}


