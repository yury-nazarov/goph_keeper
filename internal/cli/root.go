package cli

import (
	"github.com/yury-nazarov/goph_keeper/internal/models"

	"github.com/spf13/cobra"
)

// Переменные для мапинга из флагов
var (
	user   models.User
	secret models.Secret
	id     int
)

var Cmd = &cobra.Command{
	Use:   "gkc",
	Short: "Goph Keeper cli",
	Long:  `Goph Keeper command line interface`,
}

func Executor() error {
	return Cmd.Execute()
}
