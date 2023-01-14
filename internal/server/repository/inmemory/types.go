package inmemory

// Sessions - интерфейс для работы с кешем
type Sessions interface {
	//AddToken(ctx context.Context, token string, userID int) error
	//GetUserID(ctx context.Context, token string) (int, error)
	//DeleteToken(ctx context.Context, token string) error
	Close() error
}
