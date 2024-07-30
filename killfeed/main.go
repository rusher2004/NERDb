package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/killfeed/everef"
	"github.com/rusher2004/nerdb/killfeed/listener"
	"golang.org/x/sync/errgroup"
)

func main() {
	erCmd := flag.NewFlagSet("everef", flag.ExitOnError)
	erEnable := false

	lstnrCmd := flag.NewFlagSet("zkill", flag.ExitOnError)
	lstnrEnable := false

	if len(os.Args) < 2 {
		log.Fatal("must specify a subcommand")
	}

	switch os.Args[1] {
	case "everef":
		erCmd.Parse(os.Args[2:])
		erEnable = true
	case "zkill":
		lstnrCmd.Parse(os.Args[2:])
		lstnrEnable = true
	default:
		log.Fatalf("unknown subcommand: %s\n", os.Args[1])
	}

	dbURL, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatal("POSTGRES_URL is not set")
	}

	ctx := context.Background()
	cl := http.Client{Timeout: 20 * time.Second}

	pool, err := db.NewConn(ctx, dbURL)
	if err != nil {
		log.Fatalf("error creating db connection: %v", err)
	}

	done := make(chan os.Signal, 1)

	if lstnrEnable {
		log.Println("running listener")
		lstnr := listener.NewListener(&cl, *pool)
		go lstnr.Listen(ctx)
	}

	if erEnable {
		log.Println("running engine")
		g := new(errgroup.Group)
		eng := everef.NewEngine(&cl, *pool)

		g.Go(func() error {
			err := eng.RunKillmails(ctx, &cl)
			if err != nil {
				log.Printf("error running killmails: %v", err)
			}

			return err
		})

		if err := g.Wait(); err != nil {
			log.Fatalf("error running engine: %v", err)
		}

		log.Println("engine finished")
		os.Exit(0)
	}

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
