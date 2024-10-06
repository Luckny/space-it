package main

import (
	"context"
	"flag"

	"github.com/Luckny/space-it/cmd/api"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/config"
	"github.com/Luckny/space-it/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "go.uber.org/mock/gomock"
)

func main() {
	config := config.Load(".")

	addr := flag.String("addr", config.ServerAddr, "server address")
	flag.Parse()

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		util.ErrorLog.Fatal("error creating database connection", err)
	}

	if err = connPool.Ping(context.Background()); err != nil {
		util.ErrorLog.Fatal("cannot connec to database", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store, config)

	err = server.Run(*addr)
	if err != nil {
		util.ErrorLog.Fatal("cannot start the server", err)
	}
}
