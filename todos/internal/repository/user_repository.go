package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivakr07/todos/internal/models"
)

func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	var c context.Context
	var cancel context.CancelFunc
	c, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, created_at,  updated_at
	`

	err := pool.QueryRow(c, query, user.Email, user.Password).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	var c context.Context
	var cancel context.CancelFunc
	c, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, password, created_at, updated_at
		FROM users 
		WHERE email = $1
	`
	var user models.User

	err := pool.QueryRow(c, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	//here we need to return the password as our login handler need the pw to match

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func GetUserById(pool *pgxpool.Pool, id string) (*models.User, error) {
	var c context.Context
	var cancel context.CancelFunc
	c, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, created_at, updated_at
		FROM users 
		WHERE id = $1
	`
	var user models.User

	err := pool.QueryRow(c, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
