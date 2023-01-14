package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/yury-nazarov/goph_keeper/internal/server/handler"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/inmemory"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/options"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/postgres"
	"github.com/yury-nazarov/goph_keeper/internal/server/service/auth"
	"github.com/yury-nazarov/goph_keeper/internal/server/service/secret"
	"github.com/yury-nazarov/goph_keeper/pkg/application"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"

	"go.uber.org/zap"
)


// объявляем используемые зависимости и общие переменные
var (
	app      *application.Application
	log      *zap.Logger
	cfg      options.Config
	err      error
)

var (
	version = "v0.0.2"
	stage = "develop"
)

func main() {
	// Инициируем логер c нужным конфигом
	log = logger.New()

	// Читаем конфиг
	cfg, err = options.NewConfig()
	if err != nil {
		log.Fatal("can't get config", zap.String("error", err.Error()))
	}
	cfg.Version = strings.Join(serviceInfo(), " ")

	// Инициируем и запускаем приложение
	app = application.New(log, cfg, onStart, onShutdown)
	app.Run()
}

// serviceInfo получает мета информацию про сервис из переменных окружения
// 			   в противном случае используем параметры по умолчанию
func serviceInfo() (info []string) {
	if len(os.Getenv("GK_VERSION")) != 0{
		version = os.Getenv("GK_VERSION")
	}
	if len(os.Getenv("GK_STAGE")) != 0{
		stage = os.Getenv("GK_STAGE")
	}

	info = append(info, version, stage)
	return info
}

// onStart запускает проект
func onStart() {
	// Инициализируем подключение к БД
	db, err := postgres.New(log, cfg)
	if err != nil {
		log.Fatal("can't init DB storage", zap.String("error", err.Error()))
	}
	app.AddClosers(db)

	// Инициализируем подключение к кешу сессий для для хранения токенов
	sessions, err := inmemory.NewSessions(log)
	if err != nil {
		log.Fatal("can't init session cache", zap.String("error", err.Error()))
	}
	app.AddClosers(sessions)

	// Инициируем слои отвечающие за логику обработки Авторизации
	authService := auth.New(db, sessions, log)
	authController := handler.NewAuthController(authService, sessions, log)

	// Инициируем слои отвечающие за логику обработки Секретов
	secretService := secret.New(db, log)
	secretController := handler.NewSecretController(secretService, log)

	// Инициируем служебный контроллер
	// В случае с
	serviceInfo := []string{version, stage}
	msController := handler.NewMSController(serviceInfo)

	// Передаем в роутер
	r := handler.NewRouter(authController, secretController, msController)

	// Запускаем веб сервер
	go func() {
		listenAddress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		srv := http.Server{Addr: listenAddress, Handler: r}
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP Server ListenAndServe", zap.String("err", err.Error()))
		}
	}()
}

// onShutdown вызываем когда закрываем канал a.exitChannel для завершения работы приложения
func onShutdown() {
	app.ExecuteClosers()
}
