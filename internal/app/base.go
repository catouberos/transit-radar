package app

import (
	"io/fs"
	"log/slog"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/catouberos/transit-radar/internal/route"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

type App struct {
	dbPool       *pgxpool.Pool
	migrations   fs.FS
	redis        *redis.Client
	RouteService route.RouteService
}

func NewApp(dbConn *pgxpool.Pool, migrations fs.FS, redis *redis.Client) *App {
	app := &App{
		dbPool:     dbConn,
		migrations: migrations,
		redis:      redis,
	}

	return app
}

func (app *App) Init() error {
	err := app.runMigrations()
	if err != nil {
		return err
	}

	query := models.New(app.dbPool)
	app.RouteService = route.NewRouteService(query, app.redis)

	return nil
}

// runs migration if present
func (app *App) runMigrations() error {
	goose.SetBaseFS(app.migrations)

	db := stdlib.OpenDBFromPool(app.dbPool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	slog.Info("Running migrations...")

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	slog.Info("Migration completed!")

	return nil
}
