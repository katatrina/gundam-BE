// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: exchanges.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createExchange = `-- name: CreateExchange :one
INSERT INTO exchanges (id,
                       poster_id,
                       offerer_id,
                       payer_id,
                       compensation_amount,
                       status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, poster_id, offerer_id, poster_order_id, offerer_order_id, poster_from_delivery_id, poster_to_delivery_id, offerer_from_delivery_id, offerer_to_delivery_id, poster_delivery_fee, offerer_delivery_fee, poster_delivery_fee_paid, offerer_delivery_fee_paid, poster_order_expected_delivery_time, offerer_order_expected_delivery_time, poster_order_note, offerer_order_note, payer_id, compensation_amount, status, canceled_by, canceled_reason, created_at, updated_at, completed_at
`

type CreateExchangeParams struct {
	ID                 uuid.UUID      `json:"id"`
	PosterID           string         `json:"poster_id"`
	OffererID          string         `json:"offerer_id"`
	PayerID            *string        `json:"payer_id"`
	CompensationAmount *int64         `json:"compensation_amount"`
	Status             ExchangeStatus `json:"status"`
}

func (q *Queries) CreateExchange(ctx context.Context, arg CreateExchangeParams) (Exchange, error) {
	row := q.db.QueryRow(ctx, createExchange,
		arg.ID,
		arg.PosterID,
		arg.OffererID,
		arg.PayerID,
		arg.CompensationAmount,
		arg.Status,
	)
	var i Exchange
	err := row.Scan(
		&i.ID,
		&i.PosterID,
		&i.OffererID,
		&i.PosterOrderID,
		&i.OffererOrderID,
		&i.PosterFromDeliveryID,
		&i.PosterToDeliveryID,
		&i.OffererFromDeliveryID,
		&i.OffererToDeliveryID,
		&i.PosterDeliveryFee,
		&i.OffererDeliveryFee,
		&i.PosterDeliveryFeePaid,
		&i.OffererDeliveryFeePaid,
		&i.PosterOrderExpectedDeliveryTime,
		&i.OffererOrderExpectedDeliveryTime,
		&i.PosterOrderNote,
		&i.OffererOrderNote,
		&i.PayerID,
		&i.CompensationAmount,
		&i.Status,
		&i.CanceledBy,
		&i.CanceledReason,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CompletedAt,
	)
	return i, err
}

const getExchangeByID = `-- name: GetExchangeByID :one
SELECT id, poster_id, offerer_id, poster_order_id, offerer_order_id, poster_from_delivery_id, poster_to_delivery_id, offerer_from_delivery_id, offerer_to_delivery_id, poster_delivery_fee, offerer_delivery_fee, poster_delivery_fee_paid, offerer_delivery_fee_paid, poster_order_expected_delivery_time, offerer_order_expected_delivery_time, poster_order_note, offerer_order_note, payer_id, compensation_amount, status, canceled_by, canceled_reason, created_at, updated_at, completed_at
FROM exchanges
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetExchangeByID(ctx context.Context, id uuid.UUID) (Exchange, error) {
	row := q.db.QueryRow(ctx, getExchangeByID, id)
	var i Exchange
	err := row.Scan(
		&i.ID,
		&i.PosterID,
		&i.OffererID,
		&i.PosterOrderID,
		&i.OffererOrderID,
		&i.PosterFromDeliveryID,
		&i.PosterToDeliveryID,
		&i.OffererFromDeliveryID,
		&i.OffererToDeliveryID,
		&i.PosterDeliveryFee,
		&i.OffererDeliveryFee,
		&i.PosterDeliveryFeePaid,
		&i.OffererDeliveryFeePaid,
		&i.PosterOrderExpectedDeliveryTime,
		&i.OffererOrderExpectedDeliveryTime,
		&i.PosterOrderNote,
		&i.OffererOrderNote,
		&i.PayerID,
		&i.CompensationAmount,
		&i.Status,
		&i.CanceledBy,
		&i.CanceledReason,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CompletedAt,
	)
	return i, err
}

const getExchangeByOrderID = `-- name: GetExchangeByOrderID :one
SELECT id, poster_id, offerer_id, poster_order_id, offerer_order_id, poster_from_delivery_id, poster_to_delivery_id, offerer_from_delivery_id, offerer_to_delivery_id, poster_delivery_fee, offerer_delivery_fee, poster_delivery_fee_paid, offerer_delivery_fee_paid, poster_order_expected_delivery_time, offerer_order_expected_delivery_time, poster_order_note, offerer_order_note, payer_id, compensation_amount, status, canceled_by, canceled_reason, created_at, updated_at, completed_at
FROM exchanges
WHERE poster_order_id = $1
   OR offerer_order_id = $1 LIMIT 1
`

func (q *Queries) GetExchangeByOrderID(ctx context.Context, posterOrderID *uuid.UUID) (Exchange, error) {
	row := q.db.QueryRow(ctx, getExchangeByOrderID, posterOrderID)
	var i Exchange
	err := row.Scan(
		&i.ID,
		&i.PosterID,
		&i.OffererID,
		&i.PosterOrderID,
		&i.OffererOrderID,
		&i.PosterFromDeliveryID,
		&i.PosterToDeliveryID,
		&i.OffererFromDeliveryID,
		&i.OffererToDeliveryID,
		&i.PosterDeliveryFee,
		&i.OffererDeliveryFee,
		&i.PosterDeliveryFeePaid,
		&i.OffererDeliveryFeePaid,
		&i.PosterOrderExpectedDeliveryTime,
		&i.OffererOrderExpectedDeliveryTime,
		&i.PosterOrderNote,
		&i.OffererOrderNote,
		&i.PayerID,
		&i.CompensationAmount,
		&i.Status,
		&i.CanceledBy,
		&i.CanceledReason,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CompletedAt,
	)
	return i, err
}

const listUserExchanges = `-- name: ListUserExchanges :many
SELECT id, poster_id, offerer_id, poster_order_id, offerer_order_id, poster_from_delivery_id, poster_to_delivery_id, offerer_from_delivery_id, offerer_to_delivery_id, poster_delivery_fee, offerer_delivery_fee, poster_delivery_fee_paid, offerer_delivery_fee_paid, poster_order_expected_delivery_time, offerer_order_expected_delivery_time, poster_order_note, offerer_order_note, payer_id, compensation_amount, status, canceled_by, canceled_reason, created_at, updated_at, completed_at
FROM exchanges
WHERE (poster_id = $1 OR offerer_id = $1)
  AND status = coalesce($2, status)
ORDER BY created_at DESC
`

type ListUserExchangesParams struct {
	UserID string             `json:"user_id"`
	Status NullExchangeStatus `json:"status"`
}

func (q *Queries) ListUserExchanges(ctx context.Context, arg ListUserExchangesParams) ([]Exchange, error) {
	rows, err := q.db.Query(ctx, listUserExchanges, arg.UserID, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exchange{}
	for rows.Next() {
		var i Exchange
		if err := rows.Scan(
			&i.ID,
			&i.PosterID,
			&i.OffererID,
			&i.PosterOrderID,
			&i.OffererOrderID,
			&i.PosterFromDeliveryID,
			&i.PosterToDeliveryID,
			&i.OffererFromDeliveryID,
			&i.OffererToDeliveryID,
			&i.PosterDeliveryFee,
			&i.OffererDeliveryFee,
			&i.PosterDeliveryFeePaid,
			&i.OffererDeliveryFeePaid,
			&i.PosterOrderExpectedDeliveryTime,
			&i.OffererOrderExpectedDeliveryTime,
			&i.PosterOrderNote,
			&i.OffererOrderNote,
			&i.PayerID,
			&i.CompensationAmount,
			&i.Status,
			&i.CanceledBy,
			&i.CanceledReason,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CompletedAt,
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

const updateExchange = `-- name: UpdateExchange :one
UPDATE exchanges
SET poster_order_id                      = COALESCE($2, poster_order_id),
    offerer_order_id                     = COALESCE($3, offerer_order_id),
    status                               = COALESCE($4, status),

    poster_from_delivery_id              = COALESCE($5, poster_from_delivery_id),
    poster_to_delivery_id                = COALESCE($6, poster_to_delivery_id),
    offerer_from_delivery_id             = COALESCE($7, offerer_from_delivery_id),
    offerer_to_delivery_id               = COALESCE($8, offerer_to_delivery_id),

    poster_delivery_fee                  = COALESCE($9, poster_delivery_fee),
    offerer_delivery_fee                 = COALESCE($10, offerer_delivery_fee),

    poster_delivery_fee_paid             = COALESCE($11, poster_delivery_fee_paid),
    offerer_delivery_fee_paid            = COALESCE($12, offerer_delivery_fee_paid),

    poster_order_expected_delivery_time  = COALESCE($13,
                                                    poster_order_expected_delivery_time),
    offerer_order_expected_delivery_time = COALESCE($14,
                                                    offerer_order_expected_delivery_time),

    poster_order_note                    = COALESCE($15, poster_order_note),
    offerer_order_note                   = COALESCE($16, offerer_order_note),

    completed_at                         = COALESCE($17, completed_at),

    canceled_by                          = COALESCE($18, canceled_by),
    canceled_reason                      = COALESCE($19, canceled_reason),

    updated_at                           = now()
WHERE id = $1 RETURNING id, poster_id, offerer_id, poster_order_id, offerer_order_id, poster_from_delivery_id, poster_to_delivery_id, offerer_from_delivery_id, offerer_to_delivery_id, poster_delivery_fee, offerer_delivery_fee, poster_delivery_fee_paid, offerer_delivery_fee_paid, poster_order_expected_delivery_time, offerer_order_expected_delivery_time, poster_order_note, offerer_order_note, payer_id, compensation_amount, status, canceled_by, canceled_reason, created_at, updated_at, completed_at
`

type UpdateExchangeParams struct {
	ID                               uuid.UUID          `json:"id"`
	PosterOrderID                    *uuid.UUID         `json:"poster_order_id"`
	OffererOrderID                   *uuid.UUID         `json:"offerer_order_id"`
	Status                           NullExchangeStatus `json:"status"`
	PosterFromDeliveryID             *int64             `json:"poster_from_delivery_id"`
	PosterToDeliveryID               *int64             `json:"poster_to_delivery_id"`
	OffererFromDeliveryID            *int64             `json:"offerer_from_delivery_id"`
	OffererToDeliveryID              *int64             `json:"offerer_to_delivery_id"`
	PosterDeliveryFee                *int64             `json:"poster_delivery_fee"`
	OffererDeliveryFee               *int64             `json:"offerer_delivery_fee"`
	PosterDeliveryFeePaid            *bool              `json:"poster_delivery_fee_paid"`
	OffererDeliveryFeePaid           *bool              `json:"offerer_delivery_fee_paid"`
	PosterOrderExpectedDeliveryTime  *time.Time         `json:"poster_order_expected_delivery_time"`
	OffererOrderExpectedDeliveryTime *time.Time         `json:"offerer_order_expected_delivery_time"`
	PosterOrderNote                  *string            `json:"poster_order_note"`
	OffererOrderNote                 *string            `json:"offerer_order_note"`
	CompletedAt                      *time.Time         `json:"completed_at"`
	CanceledBy                       *string            `json:"canceled_by"`
	CanceledReason                   *string            `json:"canceled_reason"`
}

func (q *Queries) UpdateExchange(ctx context.Context, arg UpdateExchangeParams) (Exchange, error) {
	row := q.db.QueryRow(ctx, updateExchange,
		arg.ID,
		arg.PosterOrderID,
		arg.OffererOrderID,
		arg.Status,
		arg.PosterFromDeliveryID,
		arg.PosterToDeliveryID,
		arg.OffererFromDeliveryID,
		arg.OffererToDeliveryID,
		arg.PosterDeliveryFee,
		arg.OffererDeliveryFee,
		arg.PosterDeliveryFeePaid,
		arg.OffererDeliveryFeePaid,
		arg.PosterOrderExpectedDeliveryTime,
		arg.OffererOrderExpectedDeliveryTime,
		arg.PosterOrderNote,
		arg.OffererOrderNote,
		arg.CompletedAt,
		arg.CanceledBy,
		arg.CanceledReason,
	)
	var i Exchange
	err := row.Scan(
		&i.ID,
		&i.PosterID,
		&i.OffererID,
		&i.PosterOrderID,
		&i.OffererOrderID,
		&i.PosterFromDeliveryID,
		&i.PosterToDeliveryID,
		&i.OffererFromDeliveryID,
		&i.OffererToDeliveryID,
		&i.PosterDeliveryFee,
		&i.OffererDeliveryFee,
		&i.PosterDeliveryFeePaid,
		&i.OffererDeliveryFeePaid,
		&i.PosterOrderExpectedDeliveryTime,
		&i.OffererOrderExpectedDeliveryTime,
		&i.PosterOrderNote,
		&i.OffererOrderNote,
		&i.PayerID,
		&i.CompensationAmount,
		&i.Status,
		&i.CanceledBy,
		&i.CanceledReason,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CompletedAt,
	)
	return i, err
}
