package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/repository/client"
	"github.com/yury-nazarov/goph_keeper/internal/cli/repository/token"
	"github.com/yury-nazarov/goph_keeper/internal/cli/service/crypto"
	"github.com/yury-nazarov/goph_keeper/internal/models"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Token interface {
	Save(token string) error
	Get() (string, error)
}

type HTTPClient interface {
	Call(method string, token string, serverPath string, requestBody io.Reader) (httpStatus string, responseBody []byte, respToken string, err error)
}



type secret struct {
	token 		Token
	httpClient 	HTTPClient
}

func New() (*secret, error) {
	// Работа с токеном
	t, err := token.New()
	if err != nil {
		return nil, err
	}
	// Работа с HTTP
	c := client.New()

	s := &secret{
		token: t,
		httpClient: c,
	}
	return s, nil
}

func (s *secret) New(item models.Secret, crypto crypto.Crypto) (string, error) {
	// Шифруем секрет
	item.Data = crypto.Encrypt([]byte(item.Data))

	// Сериализуем пришедшие данные
	body, err := json.Marshal(&item)
	if err != nil {
		return  "", err
	}
	// Получаем токен
	token, err := s.token.Get()
	if err != nil {
		return "", err
	}

	// Запрос в HTTP API
	hs, _, _, err := s.httpClient.Call(http.MethodPost, token,"api/v1/secret/new", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	return hs, nil

}

func (s *secret) List(crypto crypto.Crypto) (string, []models.Secret, error){
	// Получаем токен
	token, err := s.token.Get()
	if err != nil {
		return "", nil, err
	}

	// Запрос в HTTP API
	hs, responseBody, _, err := s.httpClient.Call(http.MethodGet, token,"api/v1/secret/list", nil)
	if err != nil {
		return "", nil, err
	}

	if hs == "200 OK" {
		// Парсим JSON в слайс структур

		var items []models.Secret
		if err = json.Unmarshal(responseBody, &items); err != nil {
			return "", nil, err
		}

		// Расшифровываем содержимое
		var result []models.Secret
		for _, item := range items {
			fmt.Println(item.Data)
			item.Data, err = crypto.Decrypt(item.Data)
			fmt.Println(item.Data)
			result = append(result, item)
			if err != nil {
				log.Printf("can't decrypt message id: %d", item.ID)
			}
		}
		return hs, result, nil
	}
	return hs, nil, nil
}


func (s *secret) Get(secretID int, crypto crypto.Crypto)  (string, []models.Secret, error){
	// Получаем токен
	token, err := s.token.Get()
	if err != nil {
		return "", nil, err
	}

	// Запрос в HTTP API
	path := fmt.Sprintf("api/v1/secret/%d", secretID)
	hs, responseBody, _, err := s.httpClient.Call(http.MethodGet, token, path, nil)
	if err != nil {
		return "", nil, err
	}

	// Парсим JSON в слайс структур
	if hs == "200 OK" {
		var item models.Secret
		if err = json.Unmarshal(responseBody, &item); err != nil {
			return "", nil, err
		}
		item.Data, err = crypto.Decrypt(item.Data)
		if err != nil {
			log.Printf("can't decrypt message id: %d", item.ID)
		}

		items := []models.Secret{item}

		return hs, items, nil
	}
	return hs, nil, nil

}


func (s *secret) Delete(secretID int) (string, error){
	// Получаем токен
	token, err := s.token.Get()
	if err != nil {
		return "", err
	}

	// Запрос в HTTP API
	path := fmt.Sprintf("api/v1/secret/delete/%d", secretID)
	hs, _, _, err := s.httpClient.Call(http.MethodDelete, token, path, nil)
	if err != nil {
		return "", err
	}
	return hs, nil
}


func (s *secret) Update(item models.Secret, crypto crypto.Crypto) (string, error){
	// Получаем токен
	token, err := s.token.Get()
	if err != nil {
		return "",  err
	}

	// Запрос в HTTP API
	path := fmt.Sprintf("api/v1/secret/%d", item.ID)
	hs, responseBody, _, err := s.httpClient.Call(http.MethodGet, token, path, nil)
	if err != nil {
		return "", err
	}

	// Десериализцем в структурку
	var originSecret models.Secret
	if err = json.Unmarshal(responseBody, &originSecret); err != nil {
		return "", err
	}

	//secret - секрет измененный пользователей
	//originSecret - секрет полученый из БД.
	// Для него обновляем поля:
	//			Name, Data, Description - если они были изменены и не пустые.
	if originSecret.Name != item.Name && len(item.Name) > 0 {
		originSecret.Name = item.Name
	}
	if originSecret.Data != item.Data && len(item.Data) > 0 {
		originSecret.Data = item.Data
	}
	if originSecret.Description != item.Description && len(item.Description) > 0 {
		originSecret.Description = item.Description
	}
	// Шифруем секрет
	originSecret.Data = crypto.Encrypt([]byte(originSecret.Data))

	// Сериализуем originSecret
	body, err := json.Marshal(&originSecret)
	if err != nil {
		fmt.Println(err)
	}
	// Запрос в HTTP API для обновления данных о секрете
	hs, _, _, err = s.httpClient.Call(http.MethodPut, token, "api/v1/secret/update", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	return hs, nil
}



func (s *secret) ListOfSecrets(items []models.Secret) table.Table {
	// Формат для пользователя в терминате
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	if len(items) != 0 {
		tbl := table.New("ID", "Name", "Description", "Data")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, item := range items {
			// Расшифровываем секрет
			// TODO
			//item.Data = c.Decrypt(item.Data)
			// Добавляем в строку
			tbl.AddRow(item.ID, item.Name, item.Description, item.Data)
		}
		return tbl
	}
	return table.New()
}
