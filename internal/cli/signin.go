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
		apiServer := "http://127.0.0.1:8080/api/v1/auth/signin"

		// Готовим тело запроса в JSON
		body, err := json.Marshal(&user)
		if err != nil {
			fmt.Println(err)
		}
		// Отправляем в API на регистрацию
		resp, err := http.Post(apiServer, "application/json", bytes.NewBuffer(body))
		user.Token = resp.Header.Get("Authorization")
		defer resp.Body.Close()

		// Сохраняем данные для авторизации
		ct.AuthSave(user.Token)

		// Вывод в терминал
		fmt.Println(ct.AuthDisplayMsg(resp.Status))
	},
}

func init() {

	Cmd.AddCommand(signInCmd)
	signInCmd.Flags().StringVar(&user.Login, "login", "", "--login=username")
	signInCmd.Flags().StringVar(&user.Password, "password", "", "--password=username")
}


