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
		ct := tools.New()

		apiServer := "http://127.0.0.1:8080/api/v1/auth/signout"


		// Отправляем в API на регистрацию
		req, err := http.NewRequest(http.MethodDelete, apiServer, nil)
		if err != nil {
			fmt.Println("Connection error:", err)
		}

		ct.AuthGet()
		req.Header.Set("Authorization", ct.Token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Connection error:", err)
		}
		defer resp.Body.Close()

		// Вывод в терминал
		fmt.Println(ct.AuthDisplayMsg(resp.Status))
		// Удаляем временный файл
		ct.AuthDel()
	},
}

func init() {
	Cmd.AddCommand(signOutCmd)
}


