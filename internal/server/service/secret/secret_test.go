package secret

import (
	"context"
	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"
	_ "go.uber.org/zap"
	"reflect"
)

func (sts *StorageTestSuite) Test_secret_Create() {
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
			name: "Test_1 Создать секрет. Все поля",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
						UserID: 1,
						Name: "example",
						Data: "{\"name\": \"top secret\"}",
						Description: "example",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_2 Создать секрет обязательные поля UserID, Data",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					UserID: 1,
					Name: "example",
					Data: "{\"name\": \"top secret\"}",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_3 Не хватает обязательного поля UserID",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					Name: "example",
					Data: "{\"name\": \"top secret\"}",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_4 Не хватает обязательного поля Data",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					UserID: 1,
					Name: "example",
				},
			},
			wantErr: false,
		},
		{
			name: "Test_5 Не хватает обязательного поля Name",
			args: args{
				ctx: context.Background(),
				secret: models.Secret{
					UserID: 1,
					Data: "{\"name\": \"top secret\"}",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := &secret{
				db:  sts.TestDB,
				log: logger.New(),
			}
			if err := s.Create(tt.args.ctx, tt.args.secret); (err != nil) != tt.wantErr {
				sts.T().Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_secret_GetByID() {

	ctx := context.Background()
	ctxWithUserID := context.WithValue(ctx, "userID", 1)

	type args struct {
		ctx      context.Context
		secretID int
	}
	tests := []struct {
		name       string
		args       args
		wantSecret models.Secret
		wantErr    bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: ctxWithUserID,
				secretID: 1,
			},
			wantSecret: models.Secret{
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
			s := &secret{
				db:  sts.TestDB,
				log: logger.New(),
			}
			gotSecret, err := s.GetByID(tt.args.ctx, tt.args.secretID)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSecret, tt.wantSecret) {
				sts.T().Errorf("GetByID() gotSecret = %v, want %v", gotSecret, tt.wantSecret)
			}
		})
	}
}


func (sts *StorageTestSuite) Test_secret_List() {

	ctx := context.Background()
	ctxWithUserID := context.WithValue(ctx, "userID", 1)

	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name        string
		args        args
		wantSecrets []models.Secret
		wantErr     bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: ctxWithUserID,
				userID: 1,
			},
			wantSecrets: []models.Secret{
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
			s := &secret{
				db:  sts.TestDB,
				log: logger.New(),
			}
			gotSecrets, err := s.List(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSecrets, tt.wantSecrets) {
				sts.T().Errorf("List() gotSecrets = %v, want %v", gotSecrets, tt.wantSecrets)
			}
		})
	}
}

func (sts *StorageTestSuite) Test_secret_DeleteByID() {

	ctx := context.Background()
	ctxWithUserID := context.WithValue(ctx, "userID", 1)

	type args struct {
		ctx      context.Context
		secretID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: ctxWithUserID,
				secretID: 1,
			},
			wantErr: false,
		},

	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := &secret{
				db:  sts.TestDB,
				log: logger.New(),
			}
			if err := s.DeleteByID(tt.args.ctx, tt.args.secretID); (err != nil) != tt.wantErr {
				sts.T().Errorf("DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func (sts *StorageTestSuite) Test_secret_PutByID() {

	ctx := context.Background()
	ctxWithUserID := context.WithValue(ctx, "userID", 1)

	type args struct {
		ctx  context.Context
		item models.Secret
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_1",
			args: args{
				ctx: ctxWithUserID,
				item: models.Secret{
					ID: 1,
					UserID: 1,
					Name: "new_name_1",
					Data: "new_name_1",
					Description: "new_description_1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := &secret{
				db:  sts.TestDB,
				log: logger.New(),
			}
			if err := s.PutByID(tt.args.ctx, tt.args.item); (err != nil) != tt.wantErr {
				sts.T().Errorf("PutByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

