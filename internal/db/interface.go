package db

import (
	"context"

	"insurance_hack/internal/model"
)

type DB interface {
	GetHashedPassword(ctx context.Context, login string) (string, error)
	CreateUser(ctx context.Context, user model.User, hashedPassword string) error
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
}
