// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)

	// User
	Createuser(ctx context.Context, input RegisterUser) (output User, err error)
	GetUserByID(ctx context.Context, id int64) (output User, err error)
	GetUserByPhone(ctx context.Context, phone string) (output User, err error)
	UpdateUser(ctx context.Context, input User) (output User, err error)
	IncreaseLoginCount(ctx context.Context, id int64) (err error)
}
