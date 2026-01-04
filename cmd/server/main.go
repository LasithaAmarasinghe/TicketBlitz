package main

import (
	"database/sql"
	"fmt"
	"log"
	"ticket-blitz/internal/handlers"
	"ticket-blitz/internal/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "admin"
	password = "password123"
	dbname   = "ticketblitz"
)

func main() {
	// 1. Connect to Postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connection Pooling
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)

	// 2. Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 3. Initialize Layers (Dependency Injection)
	repo := repository.NewRepo(db, rdb)
	handler := handlers.NewTicketHandler(repo)

	// Ensure DB table exists
	if err := repo.SetupDatabase(); err != nil {
		log.Fatal("Failed to setup database:", err)
	}

	fmt.Println("🚀 TicketBlitz System Initialized (Postgres + Redis)")

	// 4. Setup Router
	r := gin.Default()
	r.POST("/reset", handler.ResetInventory)
	r.POST("/buy", handler.BuyTicket)

	// 5. Run
	if err := r.Run(":9090"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
