package auth

import (
	"context"
	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"
)


func (sts *StorageTestSuite) Test_auth_RegisterUser() {

	type args struct {
		ctx  context.Context
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1 Пользователь успешно создан",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID: 1,
					Login: "username_1",
					Password: "password_1",
					Token: "123",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_2 Имя пользователя занято",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID: 1,
					Login: "username_1",
					Password: "password_1",
					Token: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Test_3 Не заполнено имя пользователя или пароль",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID: 1,
					Login: "",
					Password: "",
					Token: "123",
				},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			a := &auth{
				log:      logger.New(),
				sessions: sts.TestSession,
				db:       sts.TestDB,
			}
			if err := a.RegisterUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				sts.T().Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_auth_UserLogIn() {

	type args struct {
		ctx  context.Context
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1 Успешный вход",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					Login: "username_1",
					Password: "password_1",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_2 Не верный логин или пароль",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					Login: "username_1",
					Password: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Test_3 Не верный логин или пароль (пустой логин или пароль)",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					Login: "",
					Password: "",
				},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			a := &auth{
				log:      logger.New(),
				sessions: sts.TestSession,
				db:       sts.TestDB,
			}
			// Регистритуем пользователя.
			_ = a.RegisterUser(tt.args.ctx, &models.User{ Login: "username_1", Password: "password_1"})
			// Проверяем вход
			if err := a.UserLogIn(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				sts.T().Errorf("UserLogIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
//
//func (sts *StorageTestSuite) Test_auth_LogOutUser() {
//
//	type args struct {
//		ctx   context.Context
//		token string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "Test_1",
//			args: args{
//				ctx: context.Background(),
//				token: "123123",
//			},
//		},
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		sts.Run(tt.name, func() {
//			a := &auth{
//				log:      logger.New(),
//				sessions: sts.TestSession,
//				db:       sts.TestDB,
//			}
//			if err := a.LogOutUser(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
//				sts.T().Errorf("LogOutUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func (sts *StorageTestSuite) Test_auth_createToken() {

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "Test_1",
			want: "123",
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			a := &auth{
				log:      logger.New(),
				sessions: sts.TestSession,
				db:       sts.TestDB,
			}
			// Функция возвращает рандомной токен, тем самым предугадать его не можем,
			// тем самым нам достаточно что это строка не 0 длины
			got, err := a.createToken()
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				sts.T().Errorf("createToken() got = %v, want %v", got, tt.want)
			}
			//if got != tt.want {
			//	sts.T().Errorf("createToken() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func (sts *StorageTestSuite) Test_auth_hashPassword() {

	type args struct {
		password string
	}
	tests := []struct {
		name   string
		args   args
		want   string
	}{
		{
			name: "Test_1",
			args: args{
				password: "123",
			},
			want: "202cb962ac59075b964b07152d234b70",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			a := &auth{
				log:      logger.New(),
				sessions: sts.TestSession,
				db:       sts.TestDB,
			}
			if got := a.hashPassword(tt.args.password); got != tt.want {
				sts.T().Errorf("hashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
