package cli

import (
	"encoding/json"
	"fmt"
	"github.com/yury-nazarov/goph_keeper/internal/models"
	"io/ioutil"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"

	"github.com/spf13/cobra"
)

var secretListCmd = &cobra.Command{
	Use:   "list",
	Short: "List of secret",
	Long:  `List of secret`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API
		resp, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/secret/list", ct.APIServer), http.MethodGet, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// JSON из набора байт
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		var secrets []models.Secret
		if err = json.Unmarshal(body, &secrets); err != nil {
			ct.Log.Warn(err.Error())
		}

		// Формат для пользователя в терминате
		// TODO:
		for k, v := range secrets {
			fmt.Printf("%d: %s\n", k, v)
		}

		// Вывод в терминал
		fmt.Println(ct.AuthDisplayMsg(resp.Status))
	},
}

func init() {
	secretCmd.AddCommand(secretListCmd)
}


