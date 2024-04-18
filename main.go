package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/ihippik/template-service/config"
	"github.com/ihippik/template-service/migrations"
	"github.com/ihippik/template-service/user"
)

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
						Action: func(ctx *cli.Context) error {
							return migrations.Up(ctx.String("conn"))
						},
					},
					{
						Name:  "down",
						Usage: "migration roll down",
						Action: func(ctx *cli.Context) error {
							return migrations.Down(ctx.String("conn"))
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/users", endpts.ListUsers)
	mux.HandleFunc("GET /v1/users/{id}", endpts.GetUser)
	mux.HandleFunc("PUT /v1/users/{id}", endpts.UpdateUser)
	mux.HandleFunc("POST /v1/users", endpts.CreateUser)
	mux.HandleFunc("DELETE /v1/users/{id}", endpts.DeleteUser)

	srv := http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           mux,
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
