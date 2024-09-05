package main

import (
	"context"

	"github.com/Luckny/space-it/cmd/api"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/config"
	"github.com/Luckny/space-it/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "go.uber.org/mock/gomock"
)

func main() {
	config := config.Load(".")

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		util.ErrorLog.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store, config)

	err = server.Run(config.ServerAddr)
	if err != nil {
		util.ErrorLog.Fatal("cannot start the server", err)
	}
}
