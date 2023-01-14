package postgres

import (
	"context"
	"reflect"

	"github.com/yury-nazarov/goph_keeper/internal/models"
)

////////// Auth

func (sts *StorageTestSuite) Test_psql_CreateUser() {
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx:      context.Background(),
				login:    "user_2",
				password: "password_2",
			},
			want:    2,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			got, err := p.CreateUser(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				sts.T().Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_psql_UserExist() {
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test_1 Пользователь существует",
			args: args{
				ctx:   context.Background(),
				login: "login_1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test_2 Пользоватлея не существует",
			args: args{
				ctx:   context.Background(),
				login: "barabashka",
			},
			want:    false,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			got, err := p.UserExist(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("UserExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				sts.T().Errorf("UserExist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_psql_UserIsValid() {
	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: context.Background(),
				user: models.User{
					Login:    "login_1",
					Password: "password_1",
				},
			},
			want:    1,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			got, err := p.UserIsValid(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("UserIsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				sts.T().Errorf("UserIsValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//////// Secrets

func (sts *StorageTestSuite) Test_psql_AddSecret() {
	type args struct {
		ctx    context.Context
		secret models.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					UserID:      4,
					Name:        "name_4",
					Data:        "data_4",
					Description: "description_4",
				},
			},
			want:    4,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			got, err := p.AddSecret(tt.args.ctx, tt.args.secret)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("AddSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				sts.T().Errorf("AddSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_psql_GetSecretByID() {
	type args struct {
		ctx    context.Context
		secret models.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    models.Secret
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					ID:     1,
					UserID: 1,
				},
			},
			want: models.Secret{
				ID:          1,
				UserID:      1,
				Name:        "name_1",
				Data:        "data_1",
				Description: "description_1",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			got, err := p.GetSecretByID(tt.args.ctx, tt.args.secret)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("GetSecretByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				sts.T().Errorf("GetSecretByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_psql_GetSecretList() {
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name           string
		args           args
		wantSecretList []models.Secret
		wantErr        bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSecretList: []models.Secret{
				{
					ID:          1,
					UserID:      1,
					Name:        "name_1",
					Data:        "data_1",
					Description: "description_1",
				},
				{
					ID:          2,
					UserID:      1,
					Name:        "name_2",
					Data:        "data_2",
					Description: "description_2",
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			gotSecretList, err := p.GetSecretList(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("GetSecretList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSecretList, tt.wantSecretList) {
				sts.T().Errorf("GetSecretList() gotSecretList = %v, want %v", gotSecretList, tt.wantSecretList)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_psql_UpdateSecretByID() {
	type args struct {
		ctx    context.Context
		secret models.Secret
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					Name:        "name_25",
					Data:        "data_25",
					Description: "description_25",
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			if err := p.UpdateSecretByID(tt.args.ctx, tt.args.secret); (err != nil) != tt.wantErr {
				sts.T().Errorf("UpdateSecretByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_psql_DeleteSecretByID() {
	type args struct {
		ctx    context.Context
		secret models.Secret
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					ID:     1,
					UserID: 1,
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorage
			if err := p.DeleteSecretByID(tt.args.ctx, tt.args.secret); (err != nil) != tt.wantErr {
				sts.T().Errorf("DeleteSecretByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
