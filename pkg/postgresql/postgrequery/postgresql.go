package postgrequery

import (
	"database/sql"
	"log"
	"time"

	"github.com/idoyudha/eshop-order/config"
	_ "github.com/lib/pq"
)

const (
	_defaultDriver       = "postgres"
	_defaultConnTimeout  = 2 * time.Second
	_defaultConnAttempts = 4 // (CPU cores × 2)
	_defaultMaxPoolSize  = 10
)

type PostgresQuery struct {
	Conn *sql.DB
}

func NewPostgres(cfg config.PostgreSQLQuery) (*PostgresQuery, error) {
	client, err := sql.Open(_defaultDriver, cfg.URL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = client.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &PostgresQuery{Conn: client}, nil
}
