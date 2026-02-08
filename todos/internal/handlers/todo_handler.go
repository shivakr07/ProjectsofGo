package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivakr07/todos/internal/repository"
)

// gin framework responsible for checking required
type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

// a struct for updated todo [as we are returning a different input]
type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// handlers
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in the context"})
			return
		}
		//why we assert type string : because our repository accepts that only
		//type comes from go itself[like any or interface{}] so we need to type annotation
		userID := userIDInterface.(string)

		var input CreateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusCreated, todo)
	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in the context"})
			return
		}
		userID := userIDInterface.(string)

		todos, err := repository.GetAllTodos(pool, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in the context"})
			return
		}
		userID := userIDInterface.(string)

		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo id"})
			return
		}

		todo, err := repository.GetTodoByID(pool, id, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, todo)
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in the context"})
			return
		}
		userID := userIDInterface.(string)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo id"})
			return
		}

		var input UpdateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Atleast one field title or completed must be provided"})
			return
		}

		existing, err := repository.GetTodoByID(pool, id, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//if todo exist
		title := existing.Title
		if input.Title != nil {
			//means user provided value of title
			title = *input.Title
		}

		completed := existing.Completed
		if input.Completed != nil {
			//means user provided value of completed
			completed = *input.Completed
		}

		todo, err := repository.UpdateTodo(pool, id, title, completed, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in the context"})
			return
		}
		userID := userIDInterface.(string)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo id"})
		}

		err = repository.DeleteTodo(pool, id, userID)
		if err != nil {
			//we match with returned error from repo
			if err.Error() == "todo with id "+idStr+" was not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo was not found"})
				return
			}

			//what if todo was found but some error occured
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo was deleted successfully!"})
	}
}
