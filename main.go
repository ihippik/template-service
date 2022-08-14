package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/ihippik/template-service/config"
	"github.com/ihippik/template-service/user"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var gitVersion = "not_specified"

// @title Swagger API ProjectName
// @version 1.0
// @description example description
// @termsOfService http://swagger.io/terms/
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host example.org
// @tag.name Template-srv
// @tag.description template service
// @BasePath /v1
func main() {
	app := &cli.App{
		Name:  "Template service",
		Usage: "template service",
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "database migration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "conn",
						Aliases:  []string{"c"},
						Usage:    "db connection",
						EnvVars:  []string{"DB_CONN"},
						Required: true,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "up",
						Usage: "migration roll up",
						Action: func(cCtx *cli.Context) error {
							db, err := sqlx.Connect("postgres", cCtx.String("conn"))
							if err != nil {
								return err
							}

							goose.SetBaseFS(embedMigrations)
							if err := goose.SetDialect("postgres"); err != nil {
								return err
							}

							if err := goose.Up(db.DB, "migrations"); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "down",
						Usage: "migration roll down",
						Action: func(cCtx *cli.Context) error {
							db, err := sqlx.Connect("postgres", cCtx.String("conn"))
							if err != nil {
								return err
							}

							goose.SetBaseFS(embedMigrations)
							if err := goose.SetDialect("postgres"); err != nil {
								return err
							}

							if err := goose.Down(db.DB, "migrations"); err != nil {
								return err
							}

							return nil
						},
					},
				},
			},
		},
		Action: func(c *cli.Context) error {
			return run(c.Context)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(mCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(mCtx, os.Interrupt)
	defer cancel()

	cfg, err := config.New(ctx)
	if err != nil {
		return fmt.Errorf("could not int config: %w", err)
	}

	logger, err := iniLogger(cfg.Log, gitVersion)
	if err != nil {
		return fmt.Errorf("could not int logger: %w", err)
	}

	db, err := initConn(cfg.DB)
	if err != nil {
		logger.Error("could`t init db connection", zap.Error(err))
		return err
	}

	svc := user.NewService(cfg, logger, user.NewRepository(db))
	endpts := user.NewEndpoint(logger, svc)

	router := mux.NewRouter()
	router.HandleFunc("/v1/users", endpts.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/v1/users/{id}", endpts.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/v1/users/{id}", endpts.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/v1/users", endpts.CreateUser).Methods(http.MethodPost)

	srv := http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
	}

	go func() {
		logger.Info("server was started", zap.String("addr", cfg.ServerAddr))

		if err := srv.ListenAndServe(); err != nil {
			logger.Error("listen & serve", zap.Error(err))
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(mCtx); err != nil {
		return err
	}

	return nil
}
