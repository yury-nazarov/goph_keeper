package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
	"net/http"
)

var signOutCmd = &cobra.Command{
	Use:   "signout",
	Short: "Logout user",
	Long:  `Logout user`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API
		resp, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/auth/signout", ct.APIServer), http.MethodDelete, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// Вывод в терминал
		fmt.Println(ct.DisplayMsg(resp.Status))

		// Удаляем временный файл
		ct.AuthDel()
	},
}

func init() {
	Cmd.AddCommand(signOutCmd)
}


