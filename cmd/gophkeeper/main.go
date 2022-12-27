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
	app      *application.Application
	db       repository.DB
	sessions repository.Sessions
	log      *zap.Logger
	cfg      options.Config
	err      error
)

func main() {
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
	// Инициализируем подключение к БД
	db, err = repository.NewPostgres(log, cfg.DB)
	if err != nil {
		log.Fatal("can't init DB storage", zap.String("error", err.Error()))
	}
	app.AddClosers(db)
	_ = db

	// Инициализируем подключение к кешу где хранятся токены
	sessions, err = repository.NewSessions(log)
	if err != nil {
		log.Fatal("can't init cache storage")
	}
	_ = sessions

	// Инициируем слой с бизнес логикой

	// Инициируем контроллер и роутер
	c := handler.NewController(db, sessions, cfg, log)
	r := handler.NewRouter(c, log)

	// Запускаем веб сервер
	go func() {
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
