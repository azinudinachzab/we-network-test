// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	InsertUser(ctx context.Context, usr User) error
	GetUserByPhone(ctx context.Context, pn string) (User, error)
	GetUserByID(ctx context.Context, id uint64) (User, error)
	UpdateUserRecord(ctx context.Context, id uint64) error
	UpdateProfileByID(ctx context.Context, id uint64, fn,pn string) error
}
