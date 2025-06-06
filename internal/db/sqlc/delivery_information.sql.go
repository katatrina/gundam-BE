// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: delivery_information.sql

package db

import (
	"context"
)

const createDeliveryInformation = `-- name: CreateDeliveryInformation :one
INSERT INTO delivery_information (user_id,
                                  full_name,
                                  phone_number,
                                  province_name,
                                  district_name,
                                  ghn_district_id,
                                  ward_name,
                                  ghn_ward_code,
                                  detail)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, user_id, full_name, phone_number, province_name, district_name, ghn_district_id, ward_name, ghn_ward_code, detail, created_at
`

type CreateDeliveryInformationParams struct {
	UserID        string `json:"user_id"`
	FullName      string `json:"full_name"`
	PhoneNumber   string `json:"phone_number"`
	ProvinceName  string `json:"province_name"`
	DistrictName  string `json:"district_name"`
	GhnDistrictID int64  `json:"ghn_district_id"`
	WardName      string `json:"ward_name"`
	GhnWardCode   string `json:"ghn_ward_code"`
	Detail        string `json:"detail"`
}

func (q *Queries) CreateDeliveryInformation(ctx context.Context, arg CreateDeliveryInformationParams) (DeliveryInformation, error) {
	row := q.db.QueryRow(ctx, createDeliveryInformation,
		arg.UserID,
		arg.FullName,
		arg.PhoneNumber,
		arg.ProvinceName,
		arg.DistrictName,
		arg.GhnDistrictID,
		arg.WardName,
		arg.GhnWardCode,
		arg.Detail,
	)
	var i DeliveryInformation
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FullName,
		&i.PhoneNumber,
		&i.ProvinceName,
		&i.DistrictName,
		&i.GhnDistrictID,
		&i.WardName,
		&i.GhnWardCode,
		&i.Detail,
		&i.CreatedAt,
	)
	return i, err
}

const getDeliveryInformation = `-- name: GetDeliveryInformation :one
SELECT id, user_id, full_name, phone_number, province_name, district_name, ghn_district_id, ward_name, ghn_ward_code, detail, created_at
FROM delivery_information
WHERE id = $1
`

func (q *Queries) GetDeliveryInformation(ctx context.Context, id int64) (DeliveryInformation, error) {
	row := q.db.QueryRow(ctx, getDeliveryInformation, id)
	var i DeliveryInformation
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FullName,
		&i.PhoneNumber,
		&i.ProvinceName,
		&i.DistrictName,
		&i.GhnDistrictID,
		&i.WardName,
		&i.GhnWardCode,
		&i.Detail,
		&i.CreatedAt,
	)
	return i, err
}
