package handler

import "github.com/dhij/ecomm/ecomm-grpc/pb"

func toPBProductReq(p ProductReq) *pb.ProductReq {
	return &pb.ProductReq{
		Id:           p.ID,
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

func toProductRes(p *pb.ProductRes) ProductRes {
	return ProductRes{
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

func toPBOrderReq(o OrderReq) *pb.OrderReq {
	return &pb.OrderReq{
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toPBOrderItems(o.Items),
	}
}

func toPBOrderItems(oi []*OrderItem) []*pb.OrderItem {
	var res []*pb.OrderItem
	for _, i := range oi {
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

func toOrderRes(o *pb.OrderRes) OrderRes {
	return OrderRes{
		ID:            o.Id,
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toOrderItems(o.Items),
	}
}

func toOrderItems(oi []*pb.OrderItem) []*OrderItem {
	var res []*OrderItem
	for _, i := range oi {
		res = append(res, &OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductID: i.ProductId,
		})
	}
	return res
}

func toPBUserReq(u UserReq) *pb.UserReq {
	return &pb.UserReq{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func toUserRes(u *pb.UserRes) UserRes {
	return UserRes{
		Name:    u.Name,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}
