package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"rpc/internal/config"
)

func main() {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/.env" // fallback для локальной разработки
	}

	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}

	log.Printf("Connecting to database: %s@%s:%d/%s",
		cfg.DbUser, cfg.DbHost, cfg.DbPort, cfg.DbName)

	ctx := context.Background()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)

	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверим подключение
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS postgres (
			id VARCHAR(36) PRIMARY KEY,
			item VARCHAR(255) NOT NULL,
			quantity INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Создаем индекс
	_, err = db.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)`)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	log.Println("Migrations completed successfully!")
}
