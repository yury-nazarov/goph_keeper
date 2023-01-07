package cli

import (
	"github.com/spf13/cobra"
	"github.com/yury-nazarov/goph_keeper/internal/models"
)

// Структура для хранения логина и токена, которые получаем во время signup, signin
// и используем в остальных методах требующих аутентификации пользователя по токену
var user models.User

var Cmd = &cobra.Command{
		Use:   "gkc",
		Short: "Goph Keeper cli",
		Long:  `Goph Keeper command line interface`,
	}

func Executor() error {
	return Cmd.Execute()
}
