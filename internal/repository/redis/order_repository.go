package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"rpc/internal/repository"
	"rpc/pkg/api/test"

	"github.com/redis/go-redis/v9"
)

type orderRepository struct {
	client *redis.Client
}

func NewOrderRepository(client *redis.Client) repository.OrderRepository {
	return &orderRepository{
		client: client,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *test.Order) error {
	key := "order:" + order.Id

	err := r.client.HSet(ctx, key,
		"id", order.Id,
		"item", order.Item,
		"quantity", order.Quantity,
	).Err()
	if err != nil {
		return fmt.Errorf("redis create: %w", err)
	}

	// TTL 10 минут
	r.client.Expire(ctx, key, 10*time.Minute)
	return nil
}

func (r *orderRepository) Get(ctx context.Context, id string) (*test.Order, error) {
	key := "order:" + id

	values, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	quantity, err := strconv.ParseInt(values["quantity"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity: %w", err)
	}

	return &test.Order{
		Id:       values["id"],
		Item:     values["item"],
		Quantity: int32(quantity),
	}, nil
}

func (r *orderRepository) Update(ctx context.Context, order *test.Order) error {
	key := "order:" + order.Id
	return r.client.Del(ctx, key).Err()
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	key := "order:" + id
	return r.client.Del(ctx, key).Err()
}

func (r *orderRepository) List(ctx context.Context) ([]*test.Order, error) {

	cached, err := r.client.Get(ctx, "orders:list").Result()
	if err != nil {
		return nil, fmt.Errorf("no list in cache: %w", err)
	}

	var orders []*test.Order
	if err := json.Unmarshal([]byte(cached), &orders); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return orders, nil
}

func (r *orderRepository) SaveList(ctx context.Context, orders []*test.Order) error {
	data, err := json.Marshal(orders)
	if err != nil {
		return fmt.Errorf("marshal orders: %w", err)
	}

	return r.client.Set(ctx, "orders:list", data, 10*time.Minute).Err()
}
