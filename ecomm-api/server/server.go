package server

import (
	"context"

	"github.com/dhij/ecomm/ecomm-api/storer"
)

type Server struct {
	storer *storer.MySQLStorer
}

func NewServer(storer *storer.MySQLStorer) *Server {
	return &Server{
		storer: storer,
	}
}

func (s *Server) CreateProduct(ctx context.Context, p *storer.Product) (*storer.Product, error) {
	return s.storer.CreateProduct(ctx, p)
}

func (s *Server) GetProduct(ctx context.Context, id int64) (*storer.Product, error) {
	return s.storer.GetProduct(ctx, id)
}

func (s *Server) ListProducts(ctx context.Context) ([]storer.Product, error) {
	return s.storer.ListProducts(ctx)
}

func (s *Server) UpdateProduct(ctx context.Context, p *storer.Product) (*storer.Product, error) {
	return s.storer.UpdateProduct(ctx, p)
}

func (s *Server) DeleteProduct(ctx context.Context, id int64) error {
	return s.storer.DeleteProduct(ctx, id)
}

func (s *Server) CreateOrder(ctx context.Context, o *storer.Order) (*storer.Order, error) {
	return s.storer.CreateOrder(ctx, o)
}

func (s *Server) GetOrder(ctx context.Context, id int64) (*storer.Order, error) {
	return s.storer.GetOrder(ctx, id)
}

func (s *Server) ListOrders(ctx context.Context) ([]storer.Order, error) {
	return s.storer.ListOrders(ctx)
}

func (s *Server) DeleteOrder(ctx context.Context, id int64) error {
	return s.storer.DeleteOrder(ctx, id)
}
