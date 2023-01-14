package postgres

// DB импортируемый интерфейс для клозера
type DB interface {
	Close() error
}
