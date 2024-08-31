package mockdb

import (
	"testing"
	"time"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/util"
	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

// randomUser generates a random db.User object
func RandomUser(t *testing.T) (db.User, string) {
	pass := util.RandomPassword()
	hash, err := util.HashPassword(pass)
	require.NoError(t, err)

	user := db.User{
		ID:        util.GenUUID(),
		Email:     util.RandomEmail(),
		Password:  hash,
		CreatedAt: pgtype.Timestamp{Time: time.Now()},
	}

	return user, pass
}

// RandomSpace generates a random db.Space object
func RandomSpace(t *testing.T, userID uuid.UUID) db.Space {

	space := db.Space{
		Name:      util.RandomSpaceName(),
		Owner:     userID,
		CreatedAt: pgtype.Timestamp{Time: time.Now()},
	}

	return space
}
