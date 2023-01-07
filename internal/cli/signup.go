package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
	"net/http"

	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create new account",
	Long:  `Create new account`,
	Run: func(cmd *cobra.Command, args []string) {
		ct := tools.New()

		// JSON для HTTP Request
		body, err := json.Marshal(&user)
		if err != nil {
			fmt.Println(err)
		}

		// Запрос в HTTP API
		resp, err := http.Post(fmt.Sprintf("%s/api/v1/auth/signup", ct.APIServer), "application/json", bytes.NewBuffer(body))
		user.Token = resp.Header.Get("Authorization")
		defer resp.Body.Close()
		// Сохраняем данные для авторизации
		ct.AuthSave(user.Token)

		// Вывод в терминал
		fmt.Println(ct.AuthDisplayMsg(resp.Status))
	},
}

func init() {
	Cmd.AddCommand(signUpCmd)
	signUpCmd.Flags().StringVarP(&user.Login, "login", "l", "", "username")
	signUpCmd.Flags().StringVarP(&user.Password, "password", "p", "", "pa$$w0rd")

	signUpCmd.MarkFlagRequired("login")
	signUpCmd.MarkFlagRequired("password")
}


