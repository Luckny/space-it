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
		ID:        util.GenUUID(),
		Name:      util.RandomSpaceName(),
		Owner:     userID,
		CreatedAt: pgtype.Timestamp{Time: time.Now()},
	}

	return space
}

// RandomPermission generates a random db.Permission object
func CreatePermission(
	t *testing.T,
	userID uuid.UUID,
	spaceID uuid.UUID,
	readPerm, writePerm, deletePerm bool,
) db.Permission {
	perm := db.Permission{
		UserID:           userID,
		SpaceID:          spaceID,
		WritePermission:  writePerm,
		ReadPermission:   readPerm,
		DeletePermission: deletePerm,
		CreatedAt:        pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:        pgtype.Timestamp{Time: time.Now()},
	}

	return perm
}

// RandomSpaceTxResult generates a random db.CreateSpaceTxResult
func RandomSpaceTxResult(t *testing.T, userId uuid.UUID) db.CreateSpaceTxResult {
	space := RandomSpace(t, userId)
	perm := CreatePermission(t, userId, space.ID, true, true, true)
	return db.CreateSpaceTxResult{
		Space:      space,
		Permission: perm,
	}
}
