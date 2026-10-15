package container

import (
	"context"
	"database/sql"

	"github.com/1337jazz/generoni/internal/adapter/config/environment"
	itemhandler "github.com/1337jazz/generoni/internal/adapter/http/handlers/item"
	"github.com/1337jazz/generoni/internal/adapter/storage/sqlite"
	"github.com/1337jazz/generoni/internal/adapter/storage/sqlite/queries"
	itemservice "github.com/1337jazz/generoni/internal/core/service/item"
)

type Container struct {
	DB          *sql.DB
	ItemHandler *itemhandler.ItemHandler
}

func New(ctx context.Context, cfg environment.Config) (*Container, error) {

	// Open the database
	db, err := sqlite.Open(ctx, cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	if err := sqlite.Migrate(cfg.DatabasePath); err != nil {
		db.Close()
		return nil, err
	}

	q := queries.New(db)
	itemRepo := sqlite.NewItemRepository(db, q)
	itemService := itemservice.NewService(itemRepo)
	itemHandler := itemhandler.NewItemHandler(itemService)

	return &Container{
		DB:          db,
		ItemHandler: itemHandler,
	}, nil
}
