package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dhij/ecomm/ecomm-grpc/pb"
	"github.com/dhij/ecomm/ecomm-grpc/storer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	storer *storer.MySQLStorer
	pb.UnimplementedEcommServer
}

func NewServer(storer *storer.MySQLStorer) *Server {
	return &Server{
		storer: storer,
	}
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	pr, err := s.storer.CreateProduct(ctx, toStorerProduct(req))
	if err != nil {
		return nil, err
	}

	return toPBProductRes(pr), nil
}

func (s *Server) GetProduct(ctx context.Context, p *pb.ProductReq) (*pb.ProductRes, error) {
	pr, err := s.storer.GetProduct(ctx, p.GetId())
	if err != nil {
		return nil, err
	}

	return toPBProductRes(pr), nil
}

func (s *Server) ListProducts(ctx context.Context, p *pb.ProductReq) (*pb.ListProductRes, error) {
	lps, err := s.storer.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	lpr := make([]*pb.ProductRes, 0, len(lps))
	for _, lp := range lps {
		lpr = append(lpr, toPBProductRes(lp))
	}

	return &pb.ListProductRes{
		Products: lpr,
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, p *pb.ProductReq) (*pb.ProductRes, error) {
	product, err := s.storer.GetProduct(ctx, p.GetId())
	if err != nil {
		return nil, err
	}

	patchProductReq(product, p)
	pr, err := s.storer.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return toPBProductRes(pr), nil
}

func (s *Server) DeleteProduct(ctx context.Context, p *pb.ProductReq) (*pb.ProductRes, error) {
	err := s.storer.DeleteProduct(ctx, p.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.ProductRes{}, nil
}

func (s *Server) CreateOrder(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	order, err := s.storer.CreateOrder(ctx, toStorerOrder(o))
	if err != nil {
		return nil, err
	}
	order.Status = storer.Pending

	_, err = s.storer.EnqueueNotificationEvent(ctx, &storer.NotificationEvent{
		UserEmail:   o.GetUserEmail(),
		OrderStatus: order.Status,
		OrderID:     order.ID,
		Attempts:    0,
	})
	if err != nil {
		return nil, err
	}

	return toPBOrderRes(order), nil
}

func (s *Server) GetOrder(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	order, err := s.storer.GetOrder(ctx, o.GetUserId())
	if err != nil {
		return nil, err
	}

	return toPBOrderRes(order), nil
}

func (s *Server) ListOrders(ctx context.Context, o *pb.OrderReq) (*pb.ListOrderRes, error) {
	orders, err := s.storer.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	lor := make([]*pb.OrderRes, 0, len(orders))
	for _, order := range orders {
		lor = append(lor, toPBOrderRes(order))
	}

	return &pb.ListOrderRes{
		Orders: lor,
	}, nil
}

func (s *Server) UpdateOrderStatus(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	// vadliate the order req
	order, err := s.storer.GetOrderStatusByID(ctx, o.GetId())
	if err != nil {
		return nil, err
	}

	if o.GetUserId() != order.UserID {
		return nil, fmt.Errorf("order %d does not belong to user %d", o.GetId(), o.GetUserId())
	}

	sOrderStatus := storer.OrderStatus(strings.ToLower(o.GetStatus().String()))
	if sOrderStatus == order.Status {
		return nil, fmt.Errorf("order status is already %s", order.Status)
	}

	order.Status = sOrderStatus
	order.UpdatedAt = toTimePtr(time.Now())
	or, err := s.storer.UpdateOrderStatus(ctx, order)
	if err != nil {
		return nil, err
	}

	// enqueue notification event
	_, err = s.storer.EnqueueNotificationEvent(ctx, &storer.NotificationEvent{
		UserEmail:   o.GetUserEmail(),
		OrderStatus: order.Status,
		OrderID:     order.ID,
		Attempts:    0,
	})
	if err != nil {
		return nil, err
	}

	return toPBOrderRes(or), nil
}

func (s *Server) DeleteOrder(ctx context.Context, o *pb.OrderReq) (*pb.OrderRes, error) {
	err := s.storer.DeleteOrder(ctx, o.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.OrderRes{}, nil
}

func (s *Server) CreateUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	user, err := s.storer.CreateUser(ctx, toStorerUser(u))
	if err != nil {
		return nil, err
	}

	return toPBUserRes(user), nil
}

func (s *Server) GetUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	user, err := s.storer.GetUser(ctx, u.GetEmail())
	if err != nil {
		return nil, err
	}

	return toPBUserRes(user), nil
}

func (s *Server) ListUsers(ctx context.Context, u *pb.UserReq) (*pb.ListUserRes, error) {
	users, err := s.storer.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	lur := make([]*pb.UserRes, 0, len(users))
	for _, user := range users {
		lur = append(lur, toPBUserRes(user))
	}

	return &pb.ListUserRes{
		Users: lur,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	user, err := s.storer.GetUser(ctx, u.GetEmail())
	if err != nil {
		return nil, err
	}

	patchUserReq(user, u)
	ur, err := s.storer.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return toPBUserRes(ur), nil
}

func (s *Server) DeleteUser(ctx context.Context, u *pb.UserReq) (*pb.UserRes, error) {
	err := s.storer.DeleteUser(ctx, u.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.UserRes{}, nil
}

func (s *Server) CreateSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	sess, err := s.storer.CreateSession(ctx, &storer.Session{
		ID:           sr.GetId(),
		UserEmail:    sr.GetUserEmail(),
		RefreshToken: sr.GetRefreshToken(),
		IsRevoked:    sr.GetIsRevoked(),
		ExpiresAt:    sr.GetExpiresAt().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{
		Id:           sess.ID,
		UserEmail:    sess.UserEmail,
		RefreshToken: sess.RefreshToken,
		IsRevoked:    sess.IsRevoked,
		ExpiresAt:    timestamppb.New(sess.ExpiresAt),
	}, nil
}

func (s *Server) GetSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	sess, err := s.storer.GetSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{
		Id:           sess.ID,
		UserEmail:    sess.UserEmail,
		RefreshToken: sess.RefreshToken,
		IsRevoked:    sess.IsRevoked,
		ExpiresAt:    timestamppb.New(sess.ExpiresAt),
	}, nil
}

func (s *Server) RevokeSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	err := s.storer.RevokeSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{}, nil
}

func (s *Server) DeleteSession(ctx context.Context, sr *pb.SessionReq) (*pb.SessionRes, error) {
	err := s.storer.DeleteSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.SessionRes{}, nil
}

func (s *Server) ListNotificationEvents(ctx context.Context, lnr *pb.ListNotificationEventsReq) (*pb.ListNotificationEventsRes, error) {
	notificationEvents, err := s.storer.ListNotificationEvents(ctx)
	if err != nil {
		return nil, err
	}

	lners := make([]*pb.NotificationEvent, 0, len(notificationEvents))
	for _, ne := range notificationEvents {
		lners = append(lners, &pb.NotificationEvent{
			Id:          ne.ID,
			UserEmail:   ne.UserEmail,
			OrderStatus: toPBOrderStatus(ne.OrderStatus),
			OrderId:     ne.OrderID,
			StateId:     ne.StateID,
			Attempts:    ne.Attempts,
		})
	}

	return &pb.ListNotificationEventsRes{
		Events: lners,
	}, nil
}

func (s *Server) UpdateNotificationEvent(ctx context.Context, unr *pb.UpdateNotificationEventReq) (*pb.UpdateNotificationEventRes, error) {
	var responseType storer.NotificationResponseType
	switch unr.ResponseType {
	case pb.NotificationResponseType_SUCCESS:
		responseType = storer.NotificationSucess
	case pb.NotificationResponseType_FAILURE:
		responseType = storer.NotificationFailure
	default:
		return nil, fmt.Errorf("invalid response type %s", unr.ResponseType)
	}

	succeeded, err := s.storer.UpdateNotificationEvent(ctx,
		&storer.NotificationEvent{
			ID:      unr.GetId(),
			StateID: unr.GetStateId(),
		},
		&storer.NotificationState{
			Message: unr.GetMessage(),
		},
		responseType)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateNotificationEventRes{
		Succeeded: succeeded,
	}, nil
}
