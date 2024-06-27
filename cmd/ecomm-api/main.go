package main

import (
	"log"

	"github.com/dhij/ecomm/db"
	"github.com/dhij/ecomm/ecomm-api/handler"
	"github.com/dhij/ecomm/ecomm-api/server"
	"github.com/dhij/ecomm/ecomm-api/storer"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to database")

	// do something with the database
	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
