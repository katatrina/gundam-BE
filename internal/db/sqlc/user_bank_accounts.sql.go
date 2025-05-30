// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user_bank_accounts.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUserBankAccount = `-- name: CreateUserBankAccount :one
INSERT INTO user_bank_accounts (id,
                                user_id,
                                account_name,
                                account_number,
                                bank_code,
                                bank_name,
                                bank_short_name)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, account_name, account_number, bank_code, bank_name, bank_short_name, created_at, updated_at
`

type CreateUserBankAccountParams struct {
	ID            uuid.UUID `json:"id"`
	UserID        string    `json:"user_id"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	BankCode      string    `json:"bank_code"`
	BankName      string    `json:"bank_name"`
	BankShortName string    `json:"bank_short_name"`
}

func (q *Queries) CreateUserBankAccount(ctx context.Context, arg CreateUserBankAccountParams) (UserBankAccount, error) {
	row := q.db.QueryRow(ctx, createUserBankAccount,
		arg.ID,
		arg.UserID,
		arg.AccountName,
		arg.AccountNumber,
		arg.BankCode,
		arg.BankName,
		arg.BankShortName,
	)
	var i UserBankAccount
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountName,
		&i.AccountNumber,
		&i.BankCode,
		&i.BankName,
		&i.BankShortName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserBankAccount = `-- name: GetUserBankAccount :one
SELECT id, user_id, account_name, account_number, bank_code, bank_name, bank_short_name, created_at, updated_at
FROM user_bank_accounts
WHERE id = $1
  AND user_id = $2
`

type GetUserBankAccountParams struct {
	ID     uuid.UUID `json:"id"`
	UserID string    `json:"user_id"`
}

func (q *Queries) GetUserBankAccount(ctx context.Context, arg GetUserBankAccountParams) (UserBankAccount, error) {
	row := q.db.QueryRow(ctx, getUserBankAccount, arg.ID, arg.UserID)
	var i UserBankAccount
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountName,
		&i.AccountNumber,
		&i.BankCode,
		&i.BankName,
		&i.BankShortName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUserBankAccounts = `-- name: ListUserBankAccounts :many
SELECT id, user_id, account_name, account_number, bank_code, bank_name, bank_short_name, created_at, updated_at
FROM user_bank_accounts
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListUserBankAccounts(ctx context.Context, userID string) ([]UserBankAccount, error) {
	rows, err := q.db.Query(ctx, listUserBankAccounts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserBankAccount{}
	for rows.Next() {
		var i UserBankAccount
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.AccountName,
			&i.AccountNumber,
			&i.BankCode,
			&i.BankName,
			&i.BankShortName,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
