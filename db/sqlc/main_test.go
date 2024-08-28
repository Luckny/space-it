package db

import (
	"context"
	"os"
	"testing"

	"github.com/Luckny/space-it/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		util.ErrorLog.Fatal("cannot load config", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		util.ErrorLog.Fatal("cannot connect to db", err)
	}

	defer connPool.Close()

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
