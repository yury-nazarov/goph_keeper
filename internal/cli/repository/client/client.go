package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	envName = "GK_API"
	defaultAPIServer = "http://127.0.0.1:8080"
)

type client struct {
	apiServer string
}

func New() *client {
	c := &client{}
	c.init()
	return c
}

func (c *client) init() {
	// Инициализируем log-файл
	if len(os.Getenv(envName)) != 0 {
		c.apiServer = os.Getenv(envName)
	} else {
		c.apiServer = defaultAPIServer
	}
}

func (c *client) Call(method string, token string, serverPath string, requestBody io.Reader) (httpStatus string, responseBody []byte, respToken string,err error) {
	// Отправляем запрос в API
	url := fmt.Sprintf("%s/%s", c.apiServer, serverPath)

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return httpStatus, responseBody, respToken, fmt.Errorf("connection error: %s", err.Error())
	}

	if len(token) != 0 {
		req.Header.Set("Authorization", token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return httpStatus, responseBody, respToken, fmt.Errorf("connection error: %s", err.Error())
	}

	// Получаем байты из resp.Body
	responseBody, err = c.getBody(resp.Body)
	if err != nil {
		return httpStatus, responseBody, respToken, fmt.Errorf("get responseBody error: %s", err.Error())
	}

	// Аутнетификация прошла успешно
	if len(resp.Header.Get("Authorization")) != 0 {
		respToken = resp.Header.Get("Authorization")
	}

	defer resp.Body.Close()
	return resp.Status, responseBody, respToken, nil
}

// getBody получает тело запроса
func (c *client) getBody(responseBody io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}
	return body, nil
}