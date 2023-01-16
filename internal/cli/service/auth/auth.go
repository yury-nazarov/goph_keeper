package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/cli/repository/client"
	"github.com/yury-nazarov/goph_keeper/internal/cli/repository/token"

	"github.com/yury-nazarov/goph_keeper/internal/models"
)

type Token interface {
	Save(token string) error
	Get() (string, error)
}

type HTTPClient interface {
	Call(method string, token string, serverPath string, requestBody io.Reader) (httpStatus string, responseBody []byte, respToken string, err error)
}

type auth struct {
	token 		Token
	httpClient 	HTTPClient
}

func New() (*auth, error) {
	// Работа с токеном
	t, err := token.New()
	if err != nil {
		return nil, err
	}
	// Работа с HTTP
	c := client.New()

	// инициируем обхект
	a := &auth{
		token: t,
		httpClient: c,
	}
	return a, nil
}

func (a *auth) SignUp(user models.User) (string, error) {
	body, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}

	// Запрос в HTTP API
	hs, _, respToken, err := a.httpClient.Call(http.MethodPost, "","api/v1/auth/signup", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// Сохраняем токен в файл
	if err := a.token.Save(respToken); err != nil {
		return "", err
	}
	return hs, nil
}

func (a *auth) SignIn(user models.User) (string, error) {
	body, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}

	// Запрос в HTTP API
	hs, _, respToken, err := a.httpClient.Call(http.MethodPost, "","api/v1/auth/signin", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// Сохраняем токен в файл
	if err := a.token.Save(respToken); err != nil {
		return "", err
	}
	return hs, nil

}

func (a *auth) SignOut() (string, error) {
	token, err := a.token.Get()
	if err != nil {
		return "", err
	}

	// Запрос в HTTP API
	hs, _, _, err := a.httpClient.Call(http.MethodDelete, token,"api/v1/auth/signout", nil)
	if err != nil {
		return "", err
	}

	return hs, nil
}