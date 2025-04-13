package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/igorgrichanov/toDoList/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder *squirrel.StatementBuilderType
	log     *slog.Logger
}

func New(conf *config.DB, log *slog.Logger) (*Postgres, error) {
	const op = "postgres.New"
	logger := log.With(
		slog.String("op", op),
	)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Name)
	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	pgConfig.MaxConns = int32(conf.MaxConn)
	pgConfig.MaxConnLifetime = conf.ConnLifeTime
	pgConfig.MaxConnIdleTime = conf.ConnIdleTime

	var pool *pgxpool.Pool
	for conf.ConnectAttempts > 0 {
		pool, err = pgxpool.NewWithConfig(context.Background(), pgConfig)
		if err == nil {
			break
		}
		logger.Info("Trying to connect to postgres", slog.Int("attempts_left", conf.ConnectAttempts))
		time.Sleep(conf.ConnectTimeout)
		conf.ConnectAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool, 0 attempts left: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Postgres{Pool: pool, Builder: &builder, log: log}, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
