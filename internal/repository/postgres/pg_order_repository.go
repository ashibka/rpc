package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rpc/internal/repository"
	"rpc/pkg/api/test"
)

type Config struct {
	DbName          string `env:"POSTGRES_DB" env-default:"postgres"`
	DbUser          string `env:"POSTGRES_USER" env-default:"postgres"`
	DbPass          string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DbHost          string `env:"POSTGRES_HOST" env-default:"db"`
	DbPort          int    `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresVersion string `env:"POSTGRES_VERSION" env-default:"15"`
}

type orderRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewOrderRepository(db *pgxpool.Pool) repository.OrderRepository {
	return &orderRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *orderRepository) Create(ctx context.Context, order *test.Order) error {
	query, args, err := r.builder.Insert("orders").
		Columns("id", "item", "quantity").
		Values(order.Id, order.Item, order.Quantity).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
func (r *orderRepository) Get(ctx context.Context, id string) (*test.Order, error) {
	query, args, err := r.builder.Select(
		"id", "item", "quantity").
		From("orders").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var order test.Order
	err = r.db.QueryRow(ctx, query, args...).Scan(&order.Id, &order.Item, &order.Quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "order with id %s not found", id)
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *test.Order) error {
	query, args, err := r.builder.Update("orders").
		Set("item", order.Item).
		Set("quantity", order.Quantity).
		Where(squirrel.Eq{"id": order.Id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "order with id %s not found", order.Id)
	}

	return nil

}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.builder.Delete("orders").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "order with id %s not found", id)
	}

	return nil
}

func (r *orderRepository) List(ctx context.Context) ([]*test.Order, error) {
	query, args, err := r.builder.Select("id", "item", "quantity").
		From("orders").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*test.Order
	for rows.Next() {
		var order test.Order
		err := rows.Scan(&order.Id, &order.Item, &order.Quantity)
		if err != nil {
			return nil, fmt.Errorf("scanning order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return orders, nil
}
