package mysqlRepository

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	testifySuite "github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"newsapi/internal/domain/newsAgr"
)

type Suite struct {
	testifySuite.Suite
	factory       *Factory
	factoryCloser func()
	RR            struct {
		News newsAgr.Repository
	}
}

func Test_Suite(t *testing.T) {
	testifySuite.Run(t, new(Suite))
}

var (
	mysqlDSN = os.Getenv("TEST_MYSQL_DSN")
)

// newMysqlExternalFactory создает фабрику репозиториев для тестирования, реализованных с помощью подключения к mysql по DSN
func (suite *Suite) newMysqlExternalFactory(dsn string) (*Factory, func()) {
	factory, err := InitFactory(Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return factory, func() { _ = factory.Close() }
}

// newMysqlContainerFactory создает фабрику репозиториев для тестирования, реализованных с помощью postgres контейнеров
func (suite *Suite) newMysqlContainerFactory() (f *Factory, closer func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	// Поиск скриптов с миграциями
	_, b, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(b), "../../../infra/mysql/init/*.sql")
	migrations, err := filepath.Glob(migrationsDir)
	suite.Require().NoError(err)
	suite.Require().NotZero(migrations)
	container, err := mysql.Run(ctx,
		"mysql:9.4",
		mysql.WithScripts(migrations...),
		mysql.WithDatabase("test_db"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	suite.Require().NoError(err)
	dsn, err := container.ConnectionString(ctx)
	suite.Require().NoError(err)

	f, err = InitFactory(Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return f, func() {
		_ = f.Close()
		_ = container.Terminate(context.Background())
	}
}

// SetupTest выполняется перед каждым тестом, связанным с testifySuite
func (suite *Suite) SetupTest() {
	// Инициализация фабрики репозиториев
	if mysqlDSN != "" {
		suite.factory, suite.factoryCloser = suite.newMysqlExternalFactory(mysqlDSN)
	} else {
		suite.factory, suite.factoryCloser = suite.newMysqlContainerFactory()
	}

	// Инициализация репозиториев
	suite.RR.News = suite.factory.NewNewsRepository()
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *Suite) TearDownSubTest() {
	err := suite.factory.Cleanup()
	suite.Require().NoError(err)
}

// TearDownTest выполняется после каждого теста, связанного с suite
func (suite *Suite) TearDownTest() {
	suite.factoryCloser()
}
