// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateRequestLog(ctx context.Context, arg CreateRequestLogParams) (uuid.UUID, error)
	CreateResponseLog(ctx context.Context, arg CreateResponseLogParams) error
	CreateSpace(ctx context.Context, arg CreateSpaceParams) (Space, error)
	DeleteSpace(ctx context.Context, id uuid.UUID) error
	GetSpaceByID(ctx context.Context, id uuid.UUID) (Space, error)
	GetSpaceByName(ctx context.Context, name string) (Space, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListSpaces(ctx context.Context, arg ListSpacesParams) ([]Space, error)
	RegisterUser(ctx context.Context, arg RegisterUserParams) (User, error)
	UpdateSpace(ctx context.Context, arg UpdateSpaceParams) (Space, error)
}

var _ Querier = (*Queries)(nil)
