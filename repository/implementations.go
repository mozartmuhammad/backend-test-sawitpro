package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"
)

var ErrDuplicateData = errors.New("duplicate data")

func (r *Repository) Createuser(ctx context.Context, input RegisterUser) (output User, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return output, err
	}

	query, err := tx.PrepareContext(ctx, InsertUserQuery)
	if err != nil {
		return output, err
	}
	defer query.Close()

	var id int64
	err = query.QueryRowContext(ctx,
		input.Phone,
		input.Name,
		input.Password,
	).Scan(&id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				return output, ErrDuplicateData
			}
		}
		return output, err
	}

	err = tx.Commit()
	if err != nil {
		return output, err
	}

	output.ID = id
	return
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (output User, err error) {
	err = r.Db.QueryRowContext(ctx, GetUserByIDQuery, id).Scan(
		&output.ID,
		&output.Phone,
		&output.Name,
		&output.Password,
		&output.CreatedAt,
		&output.UpdateAt,
	)
	return
}

func (r *Repository) UpdateUser(ctx context.Context, input User) (output User, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return output, err
	}

	query, err := tx.PrepareContext(ctx, UpdateUserQuery)
	if err != nil {
		return output, err
	}
	defer query.Close()

	_, err = query.ExecContext(ctx,
		input.Phone,
		input.Name,
		// WHERE
		input.ID,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code.Name() == "unique_violation" {
				return output, ErrDuplicateData
			}
		}
		return output, err
	}

	err = tx.Commit()
	if err != nil {
		return output, err
	}

	return input, nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (output User, err error) {
	err = r.Db.QueryRowContext(ctx, GetUserByPhoneQuery, phone).Scan(
		&output.ID,
		&output.Phone,
		&output.Name,
		&output.Password,
		&output.CreatedAt,
		&output.UpdateAt,
	)
	return
}

func (r *Repository) IncreaseLoginCount(ctx context.Context, id int64) (err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}

	query, err := tx.PrepareContext(ctx, UpdateLoginCount)
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.ExecContext(ctx,
		id,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
