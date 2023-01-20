package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func JSONUnmarshal(r *http.Request, anyData interface{}) error {
	bodyData, err := io.ReadAll(r.Body)
	// Проверяем что прочитали содержимое http запроса
	if err != nil || len(bodyData) == 0 {
		return fmt.Errorf("can't read http body: %s", err.Error())
	}

	// Извлекаем JSON
	err = json.Unmarshal(bodyData, &anyData)
	if err != nil {
		return fmt.Errorf("can't json unmarshal: %s", err.Error())
	}
	return nil
}

