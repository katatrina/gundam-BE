// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: auction_bids.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createAuctionBid = `-- name: CreateAuctionBid :one
INSERT INTO auction_bids (id,
                          auction_id,
                          bidder_id,
                          participant_id,
                          amount)
VALUES ($1, $2, $3, $4, $5) RETURNING id, auction_id, bidder_id, participant_id, amount, created_at
`

type CreateAuctionBidParams struct {
	ID            uuid.UUID  `json:"id"`
	AuctionID     *uuid.UUID `json:"auction_id"`
	BidderID      *string    `json:"bidder_id"`
	ParticipantID uuid.UUID  `json:"participant_id"`
	Amount        int64      `json:"amount"`
}

func (q *Queries) CreateAuctionBid(ctx context.Context, arg CreateAuctionBidParams) (AuctionBid, error) {
	row := q.db.QueryRow(ctx, createAuctionBid,
		arg.ID,
		arg.AuctionID,
		arg.BidderID,
		arg.ParticipantID,
		arg.Amount,
	)
	var i AuctionBid
	err := row.Scan(
		&i.ID,
		&i.AuctionID,
		&i.BidderID,
		&i.ParticipantID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getAuctionBidByID = `-- name: GetAuctionBidByID :one
SELECT id, auction_id, bidder_id, participant_id, amount, created_at
FROM auction_bids
WHERE id = $1
`

func (q *Queries) GetAuctionBidByID(ctx context.Context, id uuid.UUID) (AuctionBid, error) {
	row := q.db.QueryRow(ctx, getAuctionBidByID, id)
	var i AuctionBid
	err := row.Scan(
		&i.ID,
		&i.AuctionID,
		&i.BidderID,
		&i.ParticipantID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listAuctionBids = `-- name: ListAuctionBids :many
SELECT id, auction_id, bidder_id, participant_id, amount, created_at
FROM auction_bids
WHERE auction_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListAuctionBids(ctx context.Context, auctionID *uuid.UUID) ([]AuctionBid, error) {
	rows, err := q.db.Query(ctx, listAuctionBids, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AuctionBid{}
	for rows.Next() {
		var i AuctionBid
		if err := rows.Scan(
			&i.ID,
			&i.AuctionID,
			&i.BidderID,
			&i.ParticipantID,
			&i.Amount,
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

const listUserAuctionBids = `-- name: ListUserAuctionBids :many
SELECT id, auction_id, bidder_id, participant_id, amount, created_at
FROM auction_bids
WHERE bidder_id = $1
  AND auction_id = $2
ORDER BY created_at DESC
`

type ListUserAuctionBidsParams struct {
	BidderID  *string    `json:"bidder_id"`
	AuctionID *uuid.UUID `json:"auction_id"`
}

func (q *Queries) ListUserAuctionBids(ctx context.Context, arg ListUserAuctionBidsParams) ([]AuctionBid, error) {
	rows, err := q.db.Query(ctx, listUserAuctionBids, arg.BidderID, arg.AuctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AuctionBid{}
	for rows.Next() {
		var i AuctionBid
		if err := rows.Scan(
			&i.ID,
			&i.AuctionID,
			&i.BidderID,
			&i.ParticipantID,
			&i.Amount,
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
