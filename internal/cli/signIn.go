package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"

	"github.com/spf13/cobra"
)


var signInCmd = &cobra.Command{
	Use:   "signin",
	Short: "LogIn",
	Long:  `LogIn`,
	Run: func(cmd *cobra.Command, args []string) {
		ct := tools.New()

		// JSON для HTTP Request
		body, err := json.Marshal(&user)
		if err != nil {
			fmt.Println(err)
		}

		// Запрос в HTTP API
		resp, err := http.Post(fmt.Sprintf("%s/api/v1/auth/signin", ct.APIServer), "application/json", bytes.NewBuffer(body))
		user.Token = resp.Header.Get("Authorization")
		defer resp.Body.Close()
		// Сохраняем данные для авторизации
		ct.AuthSave(user.Token)

		// Вывод в терминал
		fmt.Println(ct.DisplayMsg(resp.Status))
	},
}

func init() {

	Cmd.AddCommand(signInCmd)
	signInCmd.Flags().StringVarP(&user.Login, "login", "l", "", "username")
	signInCmd.Flags().StringVarP(&user.Password, "password", "p", "", "pa$$w0rd")

	signInCmd.MarkFlagRequired("login")
	signInCmd.MarkFlagRequired("password")
}


