package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	b64 "encoding/base64"

	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"go.uber.org/zap"
)

// cliTools вспомогательная структура при работе с cli
// описывает логер с нужным конфигом, файл хранилище токена и прочие возможноые удобства для работы
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
		storage: fmt.Sprintf("%s/%s", homedir(), ".gk_cli_R2D2"),
		Log: logger.NewFile(fmt.Sprintf("%s/%s", homedir(), ".gk_cli_logs")),
		// TODO: Читать из конфигурационного файла
		APIServer: "http://127.0.0.1:8080",
	}
}

// setStorageFile - создает файл в домашней директории пользователя
func homedir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("can't get user home dir", zap.String("error", err.Error()))
	}
	return homedir
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

// DisplayMsg выводит в терминал сообщение пользователю
func (c *cliTools) DisplayMsg(httpStatus string) string {
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
func (c *cliTools) HTTPClient(apiServer string, method string, requestBody io.Reader) (httpStatus string, responseBody []byte, err error){
	// Отправляем в API на регистрацию
	req, err := http.NewRequest(method, apiServer, requestBody)
	if err != nil {
		return httpStatus, responseBody, fmt.Errorf("connection error: %s", err.Error())
	}

	c.AuthGet()
	req.Header.Set("Authorization", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	// Получаем байты из resp.Body
	responseBody, err = c.getBody(resp.Body)
	if err != nil {
		return httpStatus, responseBody, fmt.Errorf("get responseBody error: %s", err.Error())
	}
	defer resp.Body.Close()
	return resp.Status, responseBody, nil
}

// getBody получает тело запроса
func (c *cliTools) getBody(responseBody io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Encrypt шифрует данные
func (c *cliTools) Encrypt(data []byte) string {
	encData := b64.StdEncoding.EncodeToString(data)
	return encData
}

// Decrypt расшифровывает данные
func (c *cliTools) Decrypt(encData string) string {
	data, err := b64.StdEncoding.DecodeString(encData)
	if err != nil {
		c.Log.Warn("can't decrypt secret", zap.String("error", err.Error()))
	}
	return string(data)
}

// ListOfSecrets отформатированая таблица для вывода в терминал списка секретов
func (c *cliTools) ListOfSecrets(secrets []models.Secret) table.Table{
	// Формат для пользователя в терминате
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	if len(secrets) != 0 {
		tbl := table.New("ID", "Name", "Description", "Data")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, secret := range secrets {
			// Расшифровываем секрет
			secret.Data = c.Decrypt(secret.Data)
			// Добавляем в строку
			tbl.AddRow(secret.ID, secret.Name, secret.Description, secret.Data)
		}
			return tbl
	}
	return table.New()

}