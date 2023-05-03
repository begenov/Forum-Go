package webapp

import (
	"log"
	"net/http"

	"github.com/begenov/Forum-Go/config"
	"github.com/begenov/Forum-Go/internal/controller"
	"github.com/begenov/Forum-Go/internal/logger"
	"github.com/begenov/Forum-Go/internal/repository"
	"github.com/begenov/Forum-Go/internal/repository/sqlite"
	"github.com/begenov/Forum-Go/internal/service"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{cfg: cfg}
}

func (app *App) Run() error {
	l := logger.NewLogger(&log.Logger{})
	l.Info("Starting Server")
	db, err := sqlite.InitDatabase(app.cfg.Database.Driver, app.cfg.Database.Dsn)
	if err != nil {
		return err
	}
	l.Info("Database Init")
	repos := repository.NewRepositroy(db)
	l.Info("Repository Ok")
	service := service.NewService(repos)
	l.Info("Service Ok")
	controller := controller.NewController(*service)
	l.Info(app.cfg.Server.Addr)
	return http.ListenAndServe(":"+app.cfg.Server.Addr, controller.Router())
}
