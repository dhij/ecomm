package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dhij/ecomm/ecomm-grpc/pb"
	"golang.org/x/sync/semaphore"

	gomail "gopkg.in/mail.v2"
)

type AdminInfo struct {
	Email    string
	Password string
}

type Server struct {
	client    pb.EcommClient
	adminInfo *AdminInfo
}

func NewServer(client pb.EcommClient, adminInfo *AdminInfo) *Server {
	return &Server{
		client:    client,
		adminInfo: adminInfo,
	}
}

func (s *Server) Run(ctx context.Context) {
	// process notification event every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		// process notification event
		err := s.processNotificationEvents(ctx)
		if err != nil {
			fmt.Printf("failed to process notification events: %v\n", err)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) processNotificationEvents(ctx context.Context) error {
	res, err := s.client.ListNotificationEvents(ctx, &pb.ListNotificationEventsReq{})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)
	for _, ev := range res.Events {
		wg.Add(1)
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		go func(ev *pb.NotificationEvent) {
			defer sem.Release(1)
			defer wg.Done()
			err := s.sendNotification(ctx, ev)
			err = s.updateNotificationEvent(ctx, ev, err)
			if err != nil {
				fmt.Printf("processing event: %v\n", err)
			}
		}(ev)
	}

	go func() {
		wg.Wait()
	}()

	return nil
}

func (s *Server) sendNotification(ctx context.Context, ev *pb.NotificationEvent) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.adminInfo.Email)
	m.SetHeader("To", ev.UserEmail)
	m.SetHeader("Subject", "email from ecomm")
	m.SetBody("text/plain", fmt.Sprintf("Order %d is %s", ev.OrderId, strings.ToLower(ev.OrderStatus.String())))

	d := gomail.NewDialer("smtp.gmail.com", 587, s.adminInfo.Email, s.adminInfo.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *Server) updateNotificationEvent(ctx context.Context, ev *pb.NotificationEvent, err error) error {
	req := &pb.UpdateNotificationEventReq{
		Id:      ev.Id,
		StateId: ev.StateId,
		OrderId: ev.OrderId,
	}

	switch err {
	case nil:
		req.ResponseType = pb.NotificationResponseType_SUCCESS
		req.Message = "notification sent successfully"
	default:
		req.ResponseType = pb.NotificationResponseType_FAILURE
		req.Message = fmt.Sprintf("failed: %s", err)
	}

	fmt.Printf("updating event: %v\n", req)

	_, err = s.client.UpdateNotificationEvent(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update notification event: %s", err)
	}

	return nil
}
