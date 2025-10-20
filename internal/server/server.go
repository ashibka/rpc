package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rpc/pkg/api/test"
	"sync"
)

type Serv struct {
	test.UnimplementedOrderServiceServer
	orders map[string]*test.Order
	mu     sync.RWMutex
}

func NewServer() *Serv {
	return &Serv{
		orders: make(map[string]*test.Order),
	}
}

func (s *Serv) idgen() string {
	return uuid.New().String()
}

func (s *Serv) CreateOrder(ctx context.Context, req *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.idgen()

	order := &test.Order{
		Id:       id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	s.orders[id] = order

	return &test.CreateOrderResponse{
		Id: order.Id,
	}, nil
}

func (s *Serv) GetOrder(ctx context.Context, req *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ex := s.orders[req.Id]

	if !ex {
		return nil, status.Errorf(codes.NotFound, "Error! Can not get information about order with id %s, cause it does not exists", req.Id)
	}
	return &test.GetOrderResponse{
		Order: order,
	}, nil

}

func (s *Serv) UpdateOrder(ctx context.Context, req *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ex := s.orders[req.Id]
	if !ex {
		return nil, status.Errorf(codes.NotFound, "Error! Can not update order with id %s, cause it does not exists", req.Id)
	}
	order.Item = req.Item
	order.Quantity = req.Quantity
	return &test.UpdateOrderResponse{
		Order: order,
	}, nil
}

func (s *Serv) DeleteOrder(ctx context.Context, req *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ex := s.orders[req.Id]
	if !ex {
		return nil, status.Errorf(codes.NotFound, "Error! Can not delete order with id %s, cause it does not exists", req.Id)
	}
	delete(s.orders, req.Id)
	return &test.DeleteOrderResponse{Success: true}, nil
}

func (s *Serv) ListOrders(ctx context.Context, req *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]*test.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return &test.ListOrdersResponse{
		Orders: orders,
	}, nil
}
