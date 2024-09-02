package main

import (
	"context"
	"log"

	"github.com/dhij/ecomm/ecomm-grpc/pb"
	"github.com/dhij/ecomm/ecomm-notification/server"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		svcAddr    = envflag.String("GRPC_SVC_ADDR", "0.0.0.0:9091", "address where the ecomm-grpc service is listening on")
		adminEmail = envflag.String("ADMIN_EMAIL", "", "admin email")
		adminPass  = envflag.String("ADMIN_PASSWORD", "", "admin email")
	)
	envflag.Parse()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(*svcAddr, opts...)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewEcommClient(conn)
	srv := server.NewServer(client, &server.AdminInfo{
		Email:    *adminEmail,
		Password: *adminPass,
	})

	done := make(chan struct{})
	go func() {
		srv.Run(context.Background())
		done <- struct{}{}
	}()
	<-done
}
