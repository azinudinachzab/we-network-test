package repository

import (
	"context"
	"database/sql"
	"errors"
)

func (r *Repository) InsertUser(ctx context.Context, usr User) error {
	tx, err := r.Db.Begin()
	if err != nil {
		return err
	}
	defer rollback(ctx, tx)

	q := `INSERT INTO users(id, full_name, phone_number, password) VALUES($1, $2, $3, $4);`
	if _, err := tx.ExecContext(ctx, q, usr.ID, usr.FullName, usr.PhoneNumber, usr.Password); err != nil {
		return err
	}

	q = `INSERT INTO user_counts(user_id, count) VALUES($1, $2);`
	if _, err := tx.ExecContext(ctx, q, usr.ID, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, pn string) (User, error) {
	q := `SELECT id, full_name, phone_number, password FROM users WHERE phone_number=$1`
	var (
		id int64
		fn,pnr,p string
	)
	err := r.Db.QueryRowContext(ctx, q, pn).Scan(&id, &fn, &pnr, &p)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:          uint64(id),
		FullName:    fn,
		PhoneNumber: pnr,
		Password:    p,
	}, nil
}

func (r *Repository) GetUserByID(ctx context.Context, uid uint64) (User, error) {
	q := `SELECT id, full_name, phone_number, password FROM users WHERE id=$1`
	var (
		id int64
		fn,pnr,p string
	)
	err := r.Db.QueryRowContext(ctx, q, uid).Scan(&id, &fn, &pnr, &p)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:          uint64(id),
		FullName:    fn,
		PhoneNumber: pnr,
		Password:    p,
	}, nil
}

func (r *Repository) UpdateUserRecord(ctx context.Context, id uint64) error {
	q := `UPDATE user_counts SET count = count + 1 WHERE user_id = $1`
	if _, err := r.Db.ExecContext(ctx, q, id); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateProfileByID(ctx context.Context, id uint64, fn,pn string) error {
	q := `UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3`
	if _, err := r.Db.ExecContext(ctx, q, fn, pn, id); err != nil {
		return err
	}
	return nil
}

func rollback(ctx context.Context, tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		return err
	}
	return nil
}
