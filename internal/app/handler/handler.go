package handler

import (
	"github.com/yury-nazarov/goph_keeper/internal/app/repository"
	"github.com/yury-nazarov/goph_keeper/internal/options"
	"net/http"

	"go.uber.org/zap"

)

type Controller struct {
	db 		repository.Repository
	cgf 	options.Config
	log 	*zap.Logger

}

func NewController(db repository.Repository, cfg options.Config, log *zap.Logger) *Controller {
	c := &Controller{
		db: 	db,
		cgf: 	cfg,
		log: 	log,
	}
	return c
}

func (c *Controller) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
	return
}