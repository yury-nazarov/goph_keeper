package main

import (
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/app/handler"
	"github.com/yury-nazarov/goph_keeper/internal/app/repository"
	"github.com/yury-nazarov/goph_keeper/internal/options"
	"github.com/yury-nazarov/goph_keeper/pkg/application"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"

	"go.uber.org/zap"
)

// объявляем используемые зависимости и общие переменные
var (
	app 	*application.Application
	repo 	repository.Repository
	log 	*zap.Logger
	cfg 	options.Config
	err 	error
)

func main(){
	// Инициируем логер c нужным конфигом
	log = logger.New()

	// Читаем конфиг
	cfg, err = options.NewConfig()
	if err != nil {
		log.Fatal("can't get config", zap.String("error", err.Error()))
	}

	// Инициируем и запускаем приложение
	app = application.New(log, cfg, onStart, onShutdown)
	app.Run()
}


// onStart запускает проект
func onStart() {
	// Инициируем репозиторий
	repo, err = repository.New()
	if err != nil {
		log.Fatal("can't init repository", zap.String("error", err.Error()))
	}

	// Инициируем слой с бизнес логикой


	// Инициируем контроллер и роутер
	c := handler.NewController(repo, cfg, log)
	r := handler.NewRouter(c, log)

	// Запускаем веб сервер
	go func(){
		srv := http.Server{Addr: cfg.RunAddress, Handler: r}
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP Server ListenAndServe", zap.String("err", err.Error()))
		}
	}()
}

// onShutdown вызываем когда закрываем канал a.exitChannel для завершения работы приложения
func onShutdown() {
	app.ExecuteClosers()
}