package inmemory

import (
	"context"
	"testing"

	"github.com/yury-nazarov/goph_keeper/pkg/logger"
)

func Test_inmemorySessionStorage_AddToken(t *testing.T) {
	// Инициируем сесочную
	c, _ := NewSessions(logger.New())

	type args struct {
		ctx    context.Context
		token  string
		userID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Создать новую сессию:  token + userID",
			args: args{
				ctx: context.Background(),
				token: "123",
				userID: 1,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.AddToken(tt.args.ctx, tt.args.token, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("AddToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func Test_inmemorySessionStorage_GetUserID(t *testing.T) {
	// Инициируем сесочную
	c, _ := NewSessions(logger.New())
	// Добавляем тестовые данные
	c.AddToken(context.Background(), "123", 1)

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "получить userID по токену",
			args: args{
				ctx: context.Background(),
				token: "123",
			},
			want: 1,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.GetUserID(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_inmemorySessionStorage_DeleteToken(t *testing.T) {
	// Инициируем сесочную
	c, _ := NewSessions(logger.New())

	// Добавляем тестовые данные
	c.AddToken(context.Background(), "123", 1)

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Удалить сессию по токену",
			args: args{
				ctx: context.Background(),
				token: "123",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.DeleteToken(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("DeleteToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inmemorySessionStorage_Close(t *testing.T) {
	// Инициируем сесочную
	c, _ := NewSessions(logger.New())

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "завершаем работу сесочной",
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
