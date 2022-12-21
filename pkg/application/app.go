package application

import (
	"github.com/yury-nazarov/goph_keeper/internal/options"
	"go.uber.org/zap"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Application - стуктура приложения
type Application struct {
	// логер для приложения
	log 	 		*zap.Logger
	// Конфигурация приложения
	cfg 			options.Config
	// Запускает приложение
	onStart 		func()
	// Завершает работу приложения
	onShutDown  	func()
	// Канал в котором ожидаем системный вызов о завершении работы приложения
	signalChannel 	chan os.Signal
	// Закрытие канала означет завершение работы приложения
	exitChannel 	chan struct{}
	// закрываемые при завершении структуры
	closers 		[]io.Closer
}

// New инициирует новый экзепляр приложения
func New(log *zap.Logger, cfg options.Config, onStart func(), onShutDown func()) *Application{
	app := &Application{
		log: log,
		cfg: cfg,
		onStart: onStart,
		onShutDown: onShutDown,
	}
	// Канал для сигнализации graceful shutdown
	app.exitChannel = make(chan struct{})
	return app
}

// Run запускает приложение
func (a *Application) Run() {
	a.log.Info("The app was start", zap.String("listing socket", a.cfg.RunAddress))
	if a.onStart != nil {
		a.onStart()
	}
	// отслеживает сигналы для GracefulShutdown
	go a.initSignals()

	// Блокируем основную горутину, пока канал не будет закрыт.
	<-a.exitChannel
	// Завершаем работу приложения
	if a.onShutDown != nil {
		a.onShutDown()
	}
}


// initSignals перехватывает сигнал для остановки приложения
func (a *Application) initSignals() {
	// Канал ожидает системный вызов о завершении работы
	a.signalChannel = make(chan os.Signal, 1)

	// Ожидаемые системные вызовы
	signal.Notify(a.signalChannel, syscall.SIGTERM)
	signal.Notify(a.signalChannel, syscall.SIGINT)
	signal.Notify(a.signalChannel, syscall.SIGKILL)

	for {
		select {
		case s := <- a.signalChannel:
			switch s {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				// Получили сигнал, закрываем канал, подготавливаем приложение к завершению
				close(a.signalChannel)
				a.log.Warn("The app got shutdown", zap.String("syscall", s.String()))
				// Закрываем канал завершая работу приложения
				close(a.exitChannel)
				return
			}
		}
	}
}

// ExecuteClosers - завершает работу всех объектов в []app.Closers
func (a *Application) ExecuteClosers() {
	a.log.Info("Закрываем, все, что должно быть закрыто")
}

// AddClosers - добавляет io.Closer в closers
func (a *Application) AddClosers(c ...io.Closer) {
	a.closers = append(a.closers, c...)
}