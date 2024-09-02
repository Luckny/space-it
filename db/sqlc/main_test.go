package db

import (
	"context"
	"os"
	"testing"

	"github.com/Luckny/space-it/pkg/config"
	"github.com/Luckny/space-it/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), config.Envs.DBSource)
	if err != nil {
		util.ErrorLog.Fatal("cannot connect to db", err)
	}

	defer connPool.Close()

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
