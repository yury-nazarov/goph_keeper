package token

import (
	"fmt"
	"os"
)

const (
	envName = "GK_TOKEN"
	defaultFileName = ".gk_cli_r2d2"
)

// Token набор инструментов для работы клиента с токеном
type token struct {
	filePath string
}

func New() (*token, error) {
	t := &token{}
	if err := t.init(); err != nil {
		return nil, err
	}
	return t, nil
}

// init получает путь до файла с токеном либо из переменной окружения,
// 	    либо устанавливает по умолчанию в домашней директории пользователя
func (t *token) init() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Инициализируем log-файл
	if len(os.Getenv(envName)) != 0 {
		t.filePath = os.Getenv(envName)
	} else {
		t.filePath =  fmt.Sprintf("%s/%s", homedir, defaultFileName)
	}
	return nil
}



// Get вернет токен для операций которым нужна аутентификация на сервере
func (t *token) Get() (string, error) {
	item, err := os.ReadFile(t.filePath)
	if err != nil {
		return "", err

	}
	return string(item), nil
}

// Save - сохранить токен в файл.
// 			  т.к. для cli клиента нет постоянного рантайма,
// 			  сохраняем токен во временный файл до логаута.
//			  При новом логине файл перезаписывается.
func (t *token) Save(token string) error {
	file, err := os.Create(t.filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString(token); err != nil {
		return err
	}
	return nil
}


// Del - удаляет файл с токеном при логауте
func (t *token) Del() error {
	err := os.Remove(t.filePath)
	if err != nil {
		return err
	}
	return nil
}