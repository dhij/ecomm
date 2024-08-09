package server

import (
	"time"

	"github.com/dhij/ecomm/ecomm-grpc/pb"
	"github.com/dhij/ecomm/ecomm-grpc/storer"
	"github.com/dhij/ecomm/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toStorerProduct(p *pb.ProductReq) *storer.Product {
	return &storer.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toPBProductRes(p *storer.Product) *pb.ProductRes {
	res := &pb.ProductRes{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    timestamppb.New(p.CreatedAt),
	}
	if p.UpdatedAt != nil {
		res.UpdatedAt = timestamppb.New(*p.UpdatedAt)
	}

	return res
}

func patchProductReq(product *storer.Product, p *pb.ProductReq) {
	if p.Name != "" {
		product.Name = p.Name
	}
	if p.Image != "" {
		product.Image = p.Image
	}
	if p.Category != "" {
		product.Category = p.Category
	}
	if p.Description != "" {
		product.Description = p.Description
	}
	if p.Rating != 0 {
		product.Rating = p.Rating
	}
	if p.NumReviews != 0 {
		product.NumReviews = p.NumReviews
	}
	if p.Price != 0 {
		product.Price = p.Price
	}
	if p.CountInStock != 0 {
		product.CountInStock = p.CountInStock
	}
	product.UpdatedAt = toTimePtr(time.Now())
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}

func toStorerOrder(o *pb.OrderReq) *storer.Order {
	return &storer.Order{
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		UserID:        o.UserId,
		Items:         toStorerOrderItems(o.Items),
	}
}

func toStorerOrderItems(items []*pb.OrderItem) []storer.OrderItem {
	var res []storer.OrderItem
	for _, i := range items {
		res = append(res, storer.OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductID: i.ProductId,
		})
	}
	return res
}

func toPBOrderRes(o *storer.Order) *pb.OrderRes {
	res := &pb.OrderRes{
		Id:            o.ID,
		Items:         toPBOrderItems(o.Items),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     timestamppb.New(o.CreatedAt),
	}
	if o.UpdatedAt != nil {
		res.UpdatedAt = timestamppb.New(*o.UpdatedAt)
	}

	return res
}

func toPBOrderItems(items []storer.OrderItem) []*pb.OrderItem {
	var res []*pb.OrderItem
	for _, i := range items {
		res = append(res, &pb.OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductId: i.ProductID,
		})
	}
	return res
}

func toStorerUser(u *pb.UserReq) *storer.User {
	return &storer.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func toPBUserRes(u *storer.User) *pb.UserRes {
	return &pb.UserRes{
		Id:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func patchUserReq(user *storer.User, u *pb.UserReq) {
	if u.Name != "" {
		user.Name = u.Name
	}
	if u.Email != "" {
		user.Email = u.Email
	}
	if u.Password != "" {
		hashed, err := util.HashPassword(u.Password)
		if err != nil {
			panic(err)
		}
		user.Password = hashed
	}
	if u.IsAdmin {
		user.IsAdmin = u.IsAdmin
	}
	user.UpdatedAt = toTimePtr(time.Now())
}
