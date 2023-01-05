package postgres

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/yury-nazarov/goph_keeper/internal/options"
	"github.com/yury-nazarov/goph_keeper/internal/server/models"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/postgres/testhelpers"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Интекфейс TestStorage расширяет наш DB до нужных методов
type TestStorager interface {
	DB
}

// Позволяет агрегировать тесты
type StorageTestSuite struct {
	suite.Suite
	TestStorager
	container *testhelpers.TestDatabase
}

// Определяем необходимые методы для работы TestSuite

// SetupTest
func (sts *StorageTestSuite) SetupTest() {
	logger := logger.New()
	storageContainer := testhelpers.NewTestDatabase(sts.T())

	// Конфиг для подключения к БД
	opts := options.Config{
		MigrateFile: "./migrations_test",
		//MigrateTo: "02",
		DB: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
			storageContainer.Host(),
			storageContainer.Port(sts.T()),
			"postgres",
			"postgres",
			"postgres"),
	}

	store, err := New(logger, opts)
	require.NoError(sts.T(), err)

	sts.TestStorager = store
	sts.container = storageContainer
}

// TestStorageTestSuite запускает Docker
func TestStorageTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	t.Parallel()
	suite.Run(t, new(StorageTestSuite))
}

func (sts *StorageTestSuite) TearDownTest() {
	sts.container.Close(sts.T())
}


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
			name: "Crete new user",
			args: args{
				ctx: context.Background(),
				login: "user_2",
				password: "password_2",
			},
			want: 2,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p :=  sts.TestStorager
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
			name: "UserExist_1_success: Пользователь существует",
			args: args{
				ctx: context.Background(),
				login: "login_1",
			},
			want: true,
			wantErr: false,
		},
		{
			name: "UserExist_2_success: Пользоватлея не существует",
			args: args{
				ctx: context.Background(),
				login: "barabashka",
			},
			want: false,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
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
			name: "UserIsValid_1_success",
			args: args{
				ctx: context.Background(),
				user: models.User{
					Login: "login_1",
					Password: "password_1",
				},
			},
			want: 1,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
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
			name: "AddSecret #1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
						UserID: 4,
						Name: "name_4",
						Data: "data_4",
						Description: "description_4",
				},
			},
			want: 4,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
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
			name: "GetSecretByID_1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					ID: 1,
					UserID: 1,
				},
			},
			want: models.Secret{
				ID: 1,
				UserID: 1,
				Name: "name_1",
				Data: "data_1",
				Description: "description_1",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
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
			name: "GetSecretList_1",
			args: args{
				ctx: context.Background(),
				userID: 1,
			},
			wantSecretList: []models.Secret{
				{
					ID: 1,
					UserID: 1,
					Name: "name_1",
					Data: "data_1",
					Description: "description_1",
				},
				{
					ID: 2,
					UserID: 1,
					Name: "name_2",
					Data: "data_2",
					Description: "description_2",
				},
			},
			wantErr: false,

		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
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
			name: "UpdateSecretByID_1",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					Name: "name_25",
					Data: "data_25",
					Description: "description_25",
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
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
			name: "DeleteSecretByID_1_success",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					ID: 1,
					UserID: 1,
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			p := sts.TestStorager
			if err := p.DeleteSecretByID(tt.args.ctx, tt.args.secret); (err != nil) != tt.wantErr {
				sts.T().Errorf("DeleteSecretByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


