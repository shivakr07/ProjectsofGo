package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivakr07/todos/internal/config"
	"github.com/shivakr07/todos/internal/database"
	"github.com/shivakr07/todos/internal/handlers"
	"github.com/shivakr07/todos/internal/middleware"
)

func main() {

	var cfg *config.Config
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Fatal("Failed to load the database config", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect with database", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Todo API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/todos", handlers.CreateTodoHandler(pool))

	router.GET("/todos", handlers.GetAllTodosHandler(pool))

	router.GET("/todos/:id", handlers.GetTodoByIdHandler(pool))

	router.PUT("/todos/:id", handlers.UpdateTodoHandler(pool))

	router.DELETE("todos/:id", handlers.DeleteTodoHandler(pool))

	router.POST("/auth/register", handlers.CreateUserHandler(pool))

	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	//middleware test route
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectionHandler())
	//we are passing here cfg because it have our secret key

	router.Run(":" + cfg.Port)

}
