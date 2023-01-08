package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/yury-nazarov/goph_keeper/internal/cli/tools"
)

var secretGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get of secret by ID",
	Long:  `Get of secret by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// Инициируем вспомогательную структуру
		ct := tools.New()

		// Запрос в HTTP API
		resp, err := ct.HTTPClient(fmt.Sprintf("%s/api/v1/secret/%d", ct.APIServer, id), http.MethodGet, nil)
		if err != nil {
			ct.Log.Warn(err.Error())
		}
		defer resp.Body.Close()

		// JSON из набора байт
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ct.Log.Warn(err.Error())
		}

		if err = json.Unmarshal(body, &secret); err != nil {
			ct.Log.Warn(err.Error())
		}

		// Формат для пользователя в терминате
		fmt.Printf("%+v\n", secret)

		// Вывод в терминал
		fmt.Println(ct.DisplayMsg(resp.Status))
	},
}

func init() {
	secretCmd.AddCommand(secretGetCmd)
	secretGetCmd.Flags().IntVarP(&id, "id", "i",0, "id of secrets")

	secretGetCmd.MarkFlagRequired("id")
}


