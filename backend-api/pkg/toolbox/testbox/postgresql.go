package testbox

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	DefaultDBName     = "testdb"
	DefaultDBUserName = "testuser"
	DefaultDBPass     = "qwerty1234"
)

type ConnectionData struct {
	DBName   string
	UserName string
	Password string // nolint:gosec
	Host     string
	Port     int
}

func InitPostgresqlContainer(ctx context.Context, initScripts []string) (*ConnectionData, func(), error) {
	postgresCtr, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase(DefaultDBName),
		postgres.WithUsername(DefaultDBUserName),
		postgres.WithPassword(DefaultDBPass),
		postgres.WithOrderedInitScripts(initScripts...),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		return nil, nil, err
	}

	closeFunc := func() {
		err := testcontainers.TerminateContainer(postgresCtr)
		if err != nil {
			log.Info().Err(err).Msg("Error while shutting down container")
		}
	}

	host, _ := postgresCtr.Host(ctx)
	port, err := postgresCtr.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	connData := ConnectionData{
		DBName:   fmt.Sprintf("%s?sslmode=disable", DefaultDBName),
		UserName: DefaultDBUserName,
		Password: DefaultDBPass,
		Host:     host,
		Port:     int(port.Num()),
	}

	return &connData, closeFunc, nil
}
