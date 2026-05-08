package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"property-api/internal/domain"
)

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	GetAll(pageNum int, pageSize int) ([]domain.User, error)
	GetByID(id int) (domain.User, bool, error)
}

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, name, created_at, updated_at
	`

	err := r.pool.QueryRow(context.Background(), query, user.Name, user.CreatedAt, user.UpdatedAt).Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetAll(pageNum int, pageSize int) ([]domain.User, error) {
	offset := (pageNum - 1) * pageSize

	query := `
		SELECT id, name, created_at, updated_at
		FROM users
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(context.Background(), query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *PostgresUserRepository) GetByID(id int) (domain.User, bool, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, fmt.Errorf("get user by id: %w", err)
	}

	return user, true, nil
}
