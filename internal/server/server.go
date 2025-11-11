package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rpc/internal/repository"
	"rpc/pkg/api/test"
)

type Serv struct {
	test.UnimplementedOrderServiceServer
	repo repository.OrderRepository
}

func NewServer(repo repository.OrderRepository) *Serv {
	return &Serv{
		repo: repo,
	}
}

func (s *Serv) idgen() string {
	return uuid.New().String()
}

func (s *Serv) CreateOrder(ctx context.Context, req *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {

	id := s.idgen()

	order := &test.Order{
		Id:       id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &test.CreateOrderResponse{
		Id: order.Id,
	}, nil
}

func (s *Serv) GetOrder(ctx context.Context, req *test.GetOrderRequest) (*test.GetOrderResponse, error) {

	order, err := s.repo.Get(ctx, req.Id)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "order with id %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}
	return &test.GetOrderResponse{
		Order: order,
	}, nil

}

func (s *Serv) UpdateOrder(ctx context.Context, req *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	order := &test.Order{
		Id:       req.Id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}
	err := s.repo.Update(ctx, order)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "order with id %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to update order: %v", err)
	}

	return &test.UpdateOrderResponse{
		Order: order,
	}, nil
}

func (s *Serv) DeleteOrder(ctx context.Context, req *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {

	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "order with id %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete order: %v", err)
	}
	return &test.DeleteOrderResponse{Success: true}, nil
}

func (s *Serv) ListOrders(ctx context.Context, req *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {

	orders, err := s.repo.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}
	return &test.ListOrdersResponse{
		Orders: orders,
	}, nil
}
