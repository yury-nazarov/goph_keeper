package auth

import (
	"fmt"
	"testing"

	"github.com/yury-nazarov/goph_keeper/internal/server/repository/inmemory"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/options"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/postgres"
	"github.com/yury-nazarov/goph_keeper/pkg/logger"
	"github.com/yury-nazarov/goph_keeper/pkg/testhelpers"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestDB интефейс расширяет наш DB до нужных методов
type TestDB interface {
	postgres.DB
}

// TestSession интефейс расширяет наш DB до нужных методов
type TestSession interface {
	inmemory.Sessions
}

// StorageTestSuite Позволяет агрегировать тесты
type StorageTestSuite struct {
	suite.Suite
	TestDB
	TestSession
	container *testhelpers.TestDatabase
}

// Определяем необходимые методы для работы TestSuite

// SetupTest
func (sts *StorageTestSuite) SetupTest() {
	logger := logger.New()
	storageContainer := testhelpers.NewTestDatabase(sts.T())
	session, _ := inmemory.NewSessions(logger)

	// Конфиг для подключения к БД
	opts := options.Config{
		MigrateFile: "./migrations_test",
		//MigrateTo: "002",
		DB: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
			storageContainer.Host(),
			storageContainer.Port(sts.T()),
			"postgres",
			"postgres",
			"postgres"),
	}

	store, err := postgres.New(logger, opts)
	require.NoError(sts.T(), err)

	sts.TestDB = store
	sts.container = storageContainer
	sts.TestSession = session
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
