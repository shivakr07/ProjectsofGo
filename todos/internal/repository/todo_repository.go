package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivakr07/todos/internal/models"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
	//create a context for a db connection
	var ctx context.Context
	var cancel context.CancelFunc

	// we want a context which have expiration so to avoid db connection run forever
	// we schedule cancel when this function returns as we want to release the resources of context
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO todos (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`

	var todo models.Todo

	// it accepts db connection context
	// it receives query
	// it contains title and completed as we are going to give 1 and 2 to query
	// Scan : It receives the returned row from db and assign the values to particular fields of struct Todo in order
	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.Todo = []models.Todo{}

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}
