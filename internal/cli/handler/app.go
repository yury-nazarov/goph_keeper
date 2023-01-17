package handler

import (
	"log"

	"github.com/yury-nazarov/goph_keeper/internal/cli/service/auth"
	"github.com/yury-nazarov/goph_keeper/internal/cli/service/secrets"
	"github.com/yury-nazarov/goph_keeper/internal/models"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type Auth interface {
	SignUp(user models.User) (string, error)
	SignIn(user models.User) (string, error)
	SignOut() (string, error)
}

type Secret interface {
	New(item models.Secret) (string, error)
	List() (string, []models.Secret, error)
	Get(secretID int)  (string, []models.Secret, error)
	Delete(secretID int) (string, error)
	Update(item models.Secret) (string, error)
	ListOfSecrets(items []models.Secret) table.Table
}

// app 	контейнер содержащий инициализацию всего, что нужно для работы программы
// 		через него вызываем нужные методы определенных пакетов
type app struct {
	Cmd 	*cobra.Command
	User 	models.User
	Item 	models.Secret
	//Log
	Auth 	Auth
	Secret 	Secret
}

func New() *app{
	// Инициируем работу с аутентификацией
	a, err := auth.New()
	if err != nil {
		log.Fatal(err)
	}

	// Инициируем работу с секретами
	s, err := secrets.New()
	if err != nil {
		log.Fatal(err)
	}

	// Добавляем в основную обертку
	c := &app{
		Cmd: &cobra.Command{
			Use:   "gkc",
			Short: "Goph Keeper cli",
			Long:  `Goph Keeper command line interface`,
		},
		User: models.User{},
		Item: models.Secret{},
		Auth: a,
		Secret: s,
	}
	return c
}
