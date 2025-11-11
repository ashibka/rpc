package cached

import (
	"context"
	"rpc/internal/repository"
	"rpc/pkg/api/test"
)

type cachedRepository struct {
	redisRepo repository.OrderRepository
	pgRepo    repository.OrderRepository
}

func NewCachedRepository(redisRepo, pgRepo repository.OrderRepository) repository.OrderRepository {
	return &cachedRepository{
		redisRepo: redisRepo,
		pgRepo:    pgRepo,
	}
}

func (c *cachedRepository) Create(ctx context.Context, order *test.Order) error {
	err := c.pgRepo.Create(ctx, order)
	if err == nil {
		// Инвалидируем кэш при создании
		c.redisRepo.Delete(ctx, order.Id)
	}
	return err
}

func (c *cachedRepository) Get(ctx context.Context, id string) (*test.Order, error) {

	if order, err := c.redisRepo.Get(ctx, id); err == nil {
		return order, nil
	}

	order, err := c.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	c.redisRepo.Create(ctx, order)
	return order, nil
}

func (c *cachedRepository) Update(ctx context.Context, order *test.Order) error {
	err := c.pgRepo.Update(ctx, order)
	if err == nil {

		c.redisRepo.Delete(ctx, order.Id)
	}
	return err
}

func (c *cachedRepository) Delete(ctx context.Context, id string) error {
	err := c.pgRepo.Delete(ctx, id)
	if err == nil {

		c.redisRepo.Delete(ctx, id)
	}
	return err
}

func (c *cachedRepository) List(ctx context.Context) ([]*test.Order, error) {

	if orders, err := c.redisRepo.List(ctx); err == nil {
		return orders, nil
	}

	orders, err := c.pgRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	if redisRepo, ok := c.redisRepo.(interface {
		SaveList(ctx context.Context, orders []*test.Order) error
	}); ok {
		redisRepo.SaveList(ctx, orders)
	}

	return orders, nil
}
