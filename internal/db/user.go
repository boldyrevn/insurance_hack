package db

import (
	"context"
	"errors"
	"fmt"

	"insurance_hack/internal/model"
)

func (d *db) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var result model.User

	// language=PostgreSQL
	q := `
	SELECT login, first_name, last_name, age FROM users
	WHERE login = $1`

	row := d.conn.QueryRow(ctx, q, login)
	if row == nil {
		return result, fmt.Errorf("failed to find user by login")
	}

	if err := row.Scan(
		&result.Login,
		&result.FirstName,
		&result.LastName,
		&result.Age,
	); err != nil {
		return result, fmt.Errorf("failed to scan row result: %w", err)
	}

	return result, nil
}

func (d *db) CreateUser(ctx context.Context, user model.User, hashedPassword string) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}

	// language=PostgreSQL
	q1 := `
	INSERT INTO users(first_name, last_name, age, login) VALUES
	($1, $2, $3, $4)`

	// language=PostgreSQL
	q2 := `
	INSERT INTO auth(login, password) VALUES
	($1, $2)`

	_, err = tx.Exec(ctx, q1, user.FirstName, user.LastName, user.Age, user.Login)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, q2, user.Login, hashedPassword)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (d *db) GetHashedPassword(ctx context.Context, login string) (string, error) {
	// language=PostgreSQL
	q := `
	SELECT password FROM auth
	WHERE login = $1`

	rows, err := d.conn.Query(ctx, q, login)
	if err != nil {
		return "", err
	}

	if rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			return "", err
		}

		return password, nil
	}

	return "", errors.New("no user with such login")
}
