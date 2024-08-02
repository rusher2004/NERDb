package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antihax/goesi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/killfeed/everef"
	"github.com/rusher2004/nerdb/killfeed/listener"
	"github.com/rusher2004/nerdb/killfeed/updater"
	"golang.org/x/sync/errgroup"
)

func main() {
	erCmd := flag.NewFlagSet("everef", flag.ExitOnError)
	erEnable := false

	lstnrCmd := flag.NewFlagSet("zkill", flag.ExitOnError)
	lstnrEnable := false

	updateCmd := flag.NewFlagSet("updater", flag.ExitOnError)
	updaterSrc := updateCmd.String("src", "", "source")
	updateFileDir := updateCmd.String("dir", "", "dir source")
	// updateDate := updateCmd.String("date", time.Now().Format("2006-02-01"), "date source")
	// updateType := updateCmd.String("type", "", "character, corporation, or alliance")
	updateEnable := false

	if len(os.Args) < 2 {
		log.Fatal("must specify a subcommand")
	}

	switch os.Args[1] {
	case "everef":
		erCmd.Parse(os.Args[2:])
		erEnable = true
	case "updater":
		updateCmd.Parse(os.Args[2:])
		updateEnable = true
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

	pool, err := db.NewClient(ctx, dbURL)
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
			return eng.RunKillmails(ctx, &cl)
		})

		if err := g.Wait(); err != nil {
			log.Fatalf("error running engine: %v", err)
		}

		log.Println("engine finished")
		os.Exit(0)
	}

	if updateEnable {
		log.Println("running updater")

		g := new(errgroup.Group)
		eng := everef.NewEngine(&cl, *pool)

		g.Go(func() error {
			switch *updaterSrc {
			case "esi":
				ec := goesi.NewAPIClient(&cl, "nerdb - rusher2004@gmail.com - Fungus Amongus (in game)")
				u := updater.NewUpdater(*pool, ec.ESI.CharacterApi)

				for {
					if err := u.UpdateCharacters(ctx, 1000); err != nil {
						if errors.Is(err, updater.ErrNoUnnamedCharacters{}) {
							log.Println("no more unnamed characters")
							return nil
						}

						return err
					}

					log.Println("sleeping for 1 minute")
					time.Sleep(1 * time.Minute)
				}

			case "everef":
				return eng.RunPlayerUpdater(ctx, *updateFileDir)

			default:
				return errors.New("unknown source")
			}
		})

		if err := g.Wait(); err != nil {
			log.Fatalf("error running updater: %v", err)
		}

		log.Println("updater finished")
		os.Exit(0)
	}

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
