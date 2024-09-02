package main

import (
	"log"
	"net"

	"github.com/dhij/ecomm/db"
	"github.com/dhij/ecomm/ecomm-grpc/pb"
	"github.com/dhij/ecomm/ecomm-grpc/server"
	"github.com/dhij/ecomm/ecomm-grpc/storer"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
)

func main() {
	var (
		svcAddr = envflag.String("SVC_ADDR", "0.0.0.0:9091", "address where the ecomm-grpc service is listening on")
	)
	envflag.Parse()

	// instantiate db
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to database")

	// instantiate server
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)

	// register our server with the gRPC server
	grpcSrv := grpc.NewServer()
	pb.RegisterEcommServer(grpcSrv, srv)

	listener, err := net.Listen("tcp", *svcAddr)
	if err != nil {
		log.Fatalf("listener failed: %v", err)
	}

	log.Printf("server listening on %s", *svcAddr)
	err = grpcSrv.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
