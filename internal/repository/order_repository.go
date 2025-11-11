package repository

import (
	"context"
	"rpc/pkg/api/test"
)

type OrderRepository interface {
	Create(ctx context.Context, order *test.Order) error
	Get(ctx context.Context, id string) (*test.Order, error)
	Update(ctx context.Context, order *test.Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*test.Order, error)
}
