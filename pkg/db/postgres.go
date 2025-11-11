package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(config Config) (*pgxpool.Pool, error) {
	dataSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)

	pool, err := pgxpool.New(context.Background(), dataSource)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &pool, nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
