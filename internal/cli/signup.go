package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create new account",
	Long:  `Create new account`,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer := "http://127.0.0.1:8080/api/v1/auth/signup"

		// Готовим тело запроса в JSON
		body, err := json.Marshal(&user)
		if err != nil {
			fmt.Println(err)
		}
		// Отправляем в API на регистрацию
		resp, err := http.Post(apiServer, "application/json", bytes.NewBuffer(body))
		user.Token = resp.Header.Get("Authorization")
		defer resp.Body.Close()

		// Вывод в терминал
		AuthDisplayMsg(resp.Status)
	},
}

func init() {
	Cmd.AddCommand(signUpCmd)
	signUpCmd.Flags().StringVar(&user.Login, "login", "", "--login=username")
	signUpCmd.Flags().StringVar(&user.Password, "password", "", "--password=username")
}


