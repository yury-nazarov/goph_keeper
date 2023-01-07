package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yury-nazarov/goph_keeper/pkg/logger"

	"go.uber.org/zap"
)

// cliTools вспомогательная структура при работе с cli
// описывает логер с нужным конфиом, файл хранилище токена и прочие возможноые удобства для работы
type cliTools struct {
	// Доступ к токену прочитаному из файла
	Token string
	// логер
	Log *zap.Logger
	// Адрес сервера
	APIServer string
	// Путь файла куда сохраняется токен
	storage string

}

func New() *cliTools {
	return &cliTools{
		storage: setStorageFile(),
		Log: logger.New(),
		// TODO: Читать из конфигурационного файла
		APIServer: "http://127.0.0.1:8080",
	}
}

// setStorageFile - создает файл в домашней директории пользователя
func setStorageFile() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("can't get user home dir", zap.String("error", err.Error()))
	}
	return fmt.Sprintf("%s/%s", homedir, ".gkR2D2")
}

// AuthSave - сохранить токен файл.
// 			  т.к. для cli клиента нет постоянного рантайма,
// 			  сохраняем токен во временный файл до логаута.
//			  При новом логине файл перезаписывается.
func (c *cliTools) AuthSave(token string)  {
	file, err := os.Create(c.storage)
	if err != nil {
		c.Log.Fatal("can't create file", zap.String("error", err.Error()))
	}

	defer file.Close()

	fmt.Println("token:", token)
	_, err = file.WriteString(token)
	if err != nil {
		c.Log.Fatal("can't write to file", zap.String("error", err.Error()))
	}
}

// AuthGet при каждом запуске cli читает токен из файла
func (c *cliTools) AuthGet() {
	file, err := os.ReadFile(c.storage)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("can't read file: %s", c.storage), zap.String("error", err.Error()))
	}
	c.Log.Info(fmt.Sprintf("read file: %s", c.storage))
	c.Token = string(file)

}

// AuthDel - удаляет при логауте файл с токеном
func (c *cliTools) AuthDel() {
	err := os.Remove(c.storage)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("can't delete file: %s", c.storage), zap.String("error", err.Error()))
	} else {
		c.Log.Info(fmt.Sprintf("file: %s was deleted", c.storage))
	}
}

// AuthDisplayMsg выводит в терминал сообщение пользователю
func (c *cliTools) AuthDisplayMsg(httpStatus string) string {
	switch httpStatus {
	case "200 OK":
		return "Operation success"
	case "201 Created":
		return"Operation success"
	case "400 Bad Request":
		return"Request format error"
	case "401 Unauthorized":
		return"Invalid token"
	case "403 Forbidden":
		return"Incorrect login or password"
	case "409 Conflict":
		return"Login is already exist"
	case "500 Internal Server Error":
		return"Internal Server Error"
	default:
		return"Something wrong. Please try later"
	}
}

// HTTPClient метод для работы с HTTP API где нужна аутентификация по токену
func (c *cliTools) HTTPClient(apiServer string, method string, body io.Reader) (*http.Response, error){
	// Отправляем в API на регистрацию
	req, err := http.NewRequest(method, apiServer, body)
	if err != nil {
		return nil, fmt.Errorf("connection error: %s", err.Error())
	}

	c.AuthGet()
	req.Header.Set("Authorization", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection error: %s", err.Error())
	}
	return resp, nil
}

// Encrypt шифрует данные
// TODO
func (c *cliTools) Encrypt(data []byte) []byte {
	return data
}

// Decrypt расшифровывает данные
// TODO
func (c *cliTools) Decrypt(data []byte) []byte {
	return data
}