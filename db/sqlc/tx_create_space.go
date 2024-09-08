package db

import (
	"context"

	"github.com/google/uuid"
)

type CreateSpaceTxParams struct {
	Name  string    `json:"name"`
	Owner uuid.UUID `json:"owner"`
}

type CreateSpaceTxResult struct {
	Space      Space      `json:"space"`
	Permission Permission `json:"permission"`
}

// CreateSpaceTx creates a new space and gives it owner permissions
// returns created space and owner and permission object
func (store *SQLStore) CreateSpaceTx(
	ctx context.Context,
	arg CreateSpaceTxParams,
) (CreateSpaceTxResult, error) {
	var result CreateSpaceTxResult
	var txErr error

	txErr = store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Space, err = q.CreateSpace(ctx, CreateSpaceParams{
			Name:  arg.Name,
			Owner: arg.Owner,
		})
		if err != nil {
			return err
		}

		result.Permission, err = q.CreateAllPermission(ctx, CreateAllPermissionParams{
			UserID:  arg.Owner,
			SpaceID: result.Space.ID,
		})
		return err
	})

	if txErr != nil {
		return CreateSpaceTxResult{}, txErr
	}

	return result, nil
}
