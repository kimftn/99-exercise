package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"property-api/internal/domain"
)

type ListingRepository interface {
	Create(listing domain.Listing) (domain.Listing, error)
	GetAll(pageNum int, pageSize int, userID *int) ([]domain.Listing, error)
}

type PostgresListingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresListingRepository(pool *pgxpool.Pool) *PostgresListingRepository {
	return &PostgresListingRepository{pool: pool}
}

func (r *PostgresListingRepository) Create(listing domain.Listing) (domain.Listing, error) {
	query := `
		INSERT INTO listings (user_id, price, listing_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, price, listing_type, created_at, updated_at
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		listing.UserID,
		listing.Price,
		listing.ListingType,
		listing.CreatedAt,
		listing.UpdatedAt,
	).Scan(
		&listing.ID,
		&listing.UserID,
		&listing.Price,
		&listing.ListingType,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)
	if err != nil {
		return domain.Listing{}, fmt.Errorf("insert listing: %w", err)
	}

	return listing, nil
}

func (r *PostgresListingRepository) GetAll(pageNum int, pageSize int, userID *int) ([]domain.Listing, error) {
	offset := (pageNum - 1) * pageSize

	args := []any{pageSize, offset}
	query := `
		SELECT id, user_id, price, listing_type, created_at, updated_at
		FROM listings
	`

	if userID != nil {
		query += ` WHERE user_id = $3`
		args = append(args, *userID)
	}

	query += ` ORDER BY id ASC LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	listings := make([]domain.Listing, 0)
	for rows.Next() {
		var listing domain.Listing
		if err := rows.Scan(
			&listing.ID,
			&listing.UserID,
			&listing.Price,
			&listing.ListingType,
			&listing.CreatedAt,
			&listing.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	return listings, nil
}
