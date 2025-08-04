package base

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/catouberos/transit-radar/internal/queues"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type App struct {
	dbPool     *pgxpool.Pool
	migrations fs.FS

	query *models.Queries
	mux   *http.ServeMux
	queue *queues.Client
}

func NewApp(dbConn *pgxpool.Pool, migrations fs.FS, mux *http.ServeMux, queue *queues.Client) *App {
	app := &App{
		dbPool:     dbConn,
		migrations: migrations,

		mux:   mux,
		queue: queue,
	}

	return app
}

func (app *App) Serve() error {
	err := app.runMigrations()
	if err != nil {
		return err
	}

	err = app.initQueries()
	if err != nil {
		return err
	}

	return http.ListenAndServe("localhost:5000", app.mux)
}

func (app *App) Query() *models.Queries {
	return app.query
}

func (app *App) Queue() *queues.Client {
	return app.queue
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

func (app *App) initQueries() error {
	app.query = models.New(app.dbPool)

	return nil
}
