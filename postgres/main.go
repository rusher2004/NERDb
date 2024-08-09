package main

import (
	"context"
	"errors"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rusher2004/nerdb/postgres/migrate"
	"github.com/rusher2004/nerdb/postgres/ngrok"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "NERDb Database",
		Usage: "Manage and run the NERDb Postgres database",
		Commands: []*cli.Command{
			{
				Name:        "migrate",
				Usage:       "Run database migration",
				Description: "Runs database migration up to latest migration.",
				Action:      runMigration,
			},
			{
				Name:   "host",
				Usage:  "Expose Postgres container to ngrok",
				Action: grokIt,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func grokIt(ctx *cli.Context) error {
	newCtx := context.Background()

	token, ok := os.LookupEnv("NGROK_AUTHTOKEN")
	if !ok {
		return errors.New("NGROK_AUTHTOKEN is not set")
	}

	return ngrok.Listen(newCtx, token)
}

func runMigration(ctx *cli.Context) error {
	dbURL, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		return errors.New("POSTGRES_URL is not set")
	}

	return migrate.Migrate(dbURL)
}
