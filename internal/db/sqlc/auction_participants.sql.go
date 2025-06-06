// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: auction_participants.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createAuctionParticipant = `-- name: CreateAuctionParticipant :one
INSERT INTO auction_participants (id,
                                  auction_id,
                                  user_id,
                                  deposit_amount,
                                  deposit_entry_id)
VALUES ($1, $2, $3, $4, $5) RETURNING id, auction_id, user_id, deposit_amount, deposit_entry_id, is_refunded, created_at, updated_at
`

type CreateAuctionParticipantParams struct {
	ID             uuid.UUID `json:"id"`
	AuctionID      uuid.UUID `json:"auction_id"`
	UserID         string    `json:"user_id"`
	DepositAmount  int64     `json:"deposit_amount"`
	DepositEntryID int64     `json:"deposit_entry_id"`
}

func (q *Queries) CreateAuctionParticipant(ctx context.Context, arg CreateAuctionParticipantParams) (AuctionParticipant, error) {
	row := q.db.QueryRow(ctx, createAuctionParticipant,
		arg.ID,
		arg.AuctionID,
		arg.UserID,
		arg.DepositAmount,
		arg.DepositEntryID,
	)
	var i AuctionParticipant
	err := row.Scan(
		&i.ID,
		&i.AuctionID,
		&i.UserID,
		&i.DepositAmount,
		&i.DepositEntryID,
		&i.IsRefunded,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAuctionParticipantByUserID = `-- name: GetAuctionParticipantByUserID :one
SELECT id, auction_id, user_id, deposit_amount, deposit_entry_id, is_refunded, created_at, updated_at
FROM auction_participants
WHERE user_id = $1
  AND auction_id = $2
`

type GetAuctionParticipantByUserIDParams struct {
	UserID    string    `json:"user_id"`
	AuctionID uuid.UUID `json:"auction_id"`
}

func (q *Queries) GetAuctionParticipantByUserID(ctx context.Context, arg GetAuctionParticipantByUserIDParams) (AuctionParticipant, error) {
	row := q.db.QueryRow(ctx, getAuctionParticipantByUserID, arg.UserID, arg.AuctionID)
	var i AuctionParticipant
	err := row.Scan(
		&i.ID,
		&i.AuctionID,
		&i.UserID,
		&i.DepositAmount,
		&i.DepositEntryID,
		&i.IsRefunded,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAuctionParticipants = `-- name: ListAuctionParticipants :many
SELECT id, auction_id, user_id, deposit_amount, deposit_entry_id, is_refunded, created_at, updated_at
FROM auction_participants
WHERE auction_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListAuctionParticipants(ctx context.Context, auctionID uuid.UUID) ([]AuctionParticipant, error) {
	rows, err := q.db.Query(ctx, listAuctionParticipants, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AuctionParticipant{}
	for rows.Next() {
		var i AuctionParticipant
		if err := rows.Scan(
			&i.ID,
			&i.AuctionID,
			&i.UserID,
			&i.DepositAmount,
			&i.DepositEntryID,
			&i.IsRefunded,
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

const listAuctionParticipantsExcept = `-- name: ListAuctionParticipantsExcept :many
SELECT id, auction_id, user_id, deposit_amount, deposit_entry_id, is_refunded, created_at, updated_at
FROM auction_participants
WHERE auction_id = $1
  AND user_id != $2
ORDER BY created_at DESC
`

type ListAuctionParticipantsExceptParams struct {
	AuctionID uuid.UUID `json:"auction_id"`
	UserID    string    `json:"user_id"`
}

func (q *Queries) ListAuctionParticipantsExcept(ctx context.Context, arg ListAuctionParticipantsExceptParams) ([]AuctionParticipant, error) {
	rows, err := q.db.Query(ctx, listAuctionParticipantsExcept, arg.AuctionID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AuctionParticipant{}
	for rows.Next() {
		var i AuctionParticipant
		if err := rows.Scan(
			&i.ID,
			&i.AuctionID,
			&i.UserID,
			&i.DepositAmount,
			&i.DepositEntryID,
			&i.IsRefunded,
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

const updateAuctionParticipant = `-- name: UpdateAuctionParticipant :one
UPDATE auction_participants
SET is_refunded = COALESCE($2, is_refunded),
    updated_at  = now()
WHERE id = $1 RETURNING id, auction_id, user_id, deposit_amount, deposit_entry_id, is_refunded, created_at, updated_at
`

type UpdateAuctionParticipantParams struct {
	ID         uuid.UUID `json:"id"`
	IsRefunded *bool     `json:"is_refunded"`
}

func (q *Queries) UpdateAuctionParticipant(ctx context.Context, arg UpdateAuctionParticipantParams) (AuctionParticipant, error) {
	row := q.db.QueryRow(ctx, updateAuctionParticipant, arg.ID, arg.IsRefunded)
	var i AuctionParticipant
	err := row.Scan(
		&i.ID,
		&i.AuctionID,
		&i.UserID,
		&i.DepositAmount,
		&i.DepositEntryID,
		&i.IsRefunded,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
