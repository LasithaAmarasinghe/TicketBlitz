package repository

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	db  *sql.DB
	rdb *redis.Client
	ctx context.Context
}

// NewRepo creates a new repository instance
func NewRepo(db *sql.DB, rdb *redis.Client) *Repo {
	return &Repo{
		db:  db,
		rdb: rdb,
		ctx: context.Background(),
	}
}

// SetupDatabase ensures the Postgres table exists
func (r *Repo) SetupDatabase() error {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS inventory (id SERIAL PRIMARY KEY, stock INT)`)
	return err
}

// ResetInventory resets both Postgres and Redis to 100
func (r *Repo) ResetInventory() error {
	// 1. Reset Postgres
	if _, err := r.db.Exec("TRUNCATE inventory"); err != nil {
		return err
	}
	if _, err := r.db.Exec("INSERT INTO inventory (id, stock) VALUES (1, 100)"); err != nil {
		return err
	}

	// 2. Reset Redis
	return r.rdb.Set(r.ctx, "ticket_inventory", 100, 0).Err()
}

// BuyTicketAtomic attempts to decrement stock using Lua script
// Returns: true if bought, false if sold out, error if system fail
func (r *Repo) BuyTicketAtomic() (bool, error) {
	luaScript := `
		local stock = tonumber(redis.call("GET", KEYS[1]))
		if stock > 0 then
			redis.call("DECR", KEYS[1])
			return 1
		else
			return 0
		end
	`
	result, err := r.rdb.Eval(r.ctx, luaScript, []string{"ticket_inventory"}).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
