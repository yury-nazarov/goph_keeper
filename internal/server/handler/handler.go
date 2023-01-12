package handler

import (
	"net/http"
	"strings"
)

var err error

// msController - контроллер обрабатывающий служебные запросы для микросервиса
type msController struct {
	version string
}

func NewMSController(info []string) *msController {
	c := &msController{
		version: strings.Join(info, " "),
	}
	return c
}

// Version - версия и служебная информаця
func (c *msController) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.version))
	return
}