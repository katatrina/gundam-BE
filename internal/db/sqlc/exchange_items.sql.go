// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: exchange_items.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createExchangeItem = `-- name: CreateExchangeItem :one
INSERT INTO exchange_items (id,
                            exchange_id,
                            gundam_id,
                            name,
                            slug,
                            grade,
                            scale,
                            quantity,
                            weight,
                            image_url,
                            owner_id,
                            is_from_poster)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, exchange_id, gundam_id, name, slug, grade, scale, quantity, weight, image_url, owner_id, is_from_poster, created_at
`

type CreateExchangeItemParams struct {
	ID           uuid.UUID `json:"id"`
	ExchangeID   uuid.UUID `json:"exchange_id"`
	GundamID     *int64    `json:"gundam_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Grade        string    `json:"grade"`
	Scale        string    `json:"scale"`
	Quantity     int64     `json:"quantity"`
	Weight       int64     `json:"weight"`
	ImageURL     string    `json:"image_url"`
	OwnerID      *string   `json:"owner_id"`
	IsFromPoster bool      `json:"is_from_poster"`
}

func (q *Queries) CreateExchangeItem(ctx context.Context, arg CreateExchangeItemParams) (ExchangeItem, error) {
	row := q.db.QueryRow(ctx, createExchangeItem,
		arg.ID,
		arg.ExchangeID,
		arg.GundamID,
		arg.Name,
		arg.Slug,
		arg.Grade,
		arg.Scale,
		arg.Quantity,
		arg.Weight,
		arg.ImageURL,
		arg.OwnerID,
		arg.IsFromPoster,
	)
	var i ExchangeItem
	err := row.Scan(
		&i.ID,
		&i.ExchangeID,
		&i.GundamID,
		&i.Name,
		&i.Slug,
		&i.Grade,
		&i.Scale,
		&i.Quantity,
		&i.Weight,
		&i.ImageURL,
		&i.OwnerID,
		&i.IsFromPoster,
		&i.CreatedAt,
	)
	return i, err
}

const listExchangeItems = `-- name: ListExchangeItems :many
SELECT id, exchange_id, gundam_id, name, slug, grade, scale, quantity, weight, image_url, owner_id, is_from_poster, created_at
FROM exchange_items
WHERE exchange_id = $1
  AND ($2::boolean IS NULL OR is_from_poster = $2::boolean)
ORDER BY created_at DESC
`

type ListExchangeItemsParams struct {
	ExchangeID   uuid.UUID `json:"exchange_id"`
	IsFromPoster *bool     `json:"is_from_poster"`
}

func (q *Queries) ListExchangeItems(ctx context.Context, arg ListExchangeItemsParams) ([]ExchangeItem, error) {
	rows, err := q.db.Query(ctx, listExchangeItems, arg.ExchangeID, arg.IsFromPoster)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ExchangeItem{}
	for rows.Next() {
		var i ExchangeItem
		if err := rows.Scan(
			&i.ID,
			&i.ExchangeID,
			&i.GundamID,
			&i.Name,
			&i.Slug,
			&i.Grade,
			&i.Scale,
			&i.Quantity,
			&i.Weight,
			&i.ImageURL,
			&i.OwnerID,
			&i.IsFromPoster,
			&i.CreatedAt,
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
