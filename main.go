package main

import (
	"context"

	"github.com/Luckny/space-it/api"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		util.ErrorLog.Fatal("cannot load config variables", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		util.ErrorLog.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	util.InfoLog.Println("server listening on", config.ServerAddr)
	err = server.Run(config.ServerAddr)
	if err != nil {
		util.ErrorLog.Fatal("cannot start the server", err)
	}
}
