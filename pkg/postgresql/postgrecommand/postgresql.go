package postgrecommand

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/idoyudha/eshop-order/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCommand struct {
	maxConnSize  int
	minConnSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func NewPostgres(cfg config.PostgreSQLCommand) (*PostgresCommand, error) {
	pg := &PostgresCommand{
		maxConnSize:  cfg.MaxConnSize,
		minConnSize:  cfg.MinConnSize,
		connAttempts: cfg.ConnAttemps,
		connTimeout:  time.Duration(cfg.ConnTimeout),
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(pg.maxConnSize)
	poolConfig.MinConns = int32(pg.minConnSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("postgres command is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres command failed to connect, connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *PostgresCommand) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
