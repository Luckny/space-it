package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	ConnectionFailure   = "08006"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

var ErrConnectionFailure = &pgconn.PgError{
	Code: ConnectionFailure,
}

var ErrForeignKeyConstraint = &pgconn.PgError{
	Code: ForeignKeyViolation,
}
