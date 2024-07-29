package everef

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mholt/archiver/v4"
	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/killfeed/esi"
	"golang.org/x/sync/errgroup"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Engine struct {
	db db.Client
	hc HTTPClient
}

func NewEngine(hc HTTPClient, db db.Client) *Engine {
	return &Engine{db, hc}
}

func handleTarFileWithDB(cl db.Client) func(ct context.Context, f archiver.File) error {
	return func(ctx context.Context, f archiver.File) error {
		if f.IsDir() {
			log.Printf("skipping directory: %s", f.Name())
			return nil
		}

		contents, err := f.Open()
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}

		b, err := io.ReadAll(contents)
		if err != nil {
			return fmt.Errorf("error reading file contents: %w", err)
		}

		var km esi.Killmail
		if err := json.Unmarshal(b, &km); err != nil {
			return fmt.Errorf("error unmarshalling killmail: %w", err)
		}

		if err := cl.UpsertESIKillmail(ctx, km); err != nil {
			return fmt.Errorf("error upserting killmail: %w", err)
		}

		return nil
	}
}

func (e *Engine) ProcessDay(ctx context.Context, day string) error {
	// fetch killmails for the day
	// parse day into a time.Time object as YYYYMMDD
	t, err := time.Parse("20060102", day)
	if err != nil {
		return fmt.Errorf("error parsing day: %w", err)
	}

	// download from https://data.everef.net/killmails/2007/killmails-2007-12-05.tar.bz2
	uri := fmt.Sprintf("https://data.everef.net/killmails/%d/killmails-%d-%02d-%02d.tar.bz2", t.Year(), t.Year(), t.Month(), t.Day())
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	res, err := e.hc.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	bz2 := archiver.Bz2{}
	bReader, err := bz2.OpenReader(res.Body)
	if err != nil {
		return fmt.Errorf("error opening bzip2 reader: %w", err)
	}
	defer bReader.Close()

	tar := archiver.Tar{}
	if err := tar.Extract(ctx, bReader, nil, handleTarFileWithDB(e.db)); err != nil {
		return fmt.Errorf("error extracting tarball: %w", err)
	}

	return nil
}

func (e *Engine) RunKillmails(ctx context.Context, hc HTTPClient) error {
	g := new(errgroup.Group)

	totals := map[string]map[string]int{
		"db":     nil,
		"everef": nil,
	}
	g.Go(func() error {
		t, err := FetchTotals(hc)
		if err != nil {
			return fmt.Errorf("everef error: %w", err)
		}
		totals["everef"] = t

		return nil
	})

	g.Go(func() error {
		t, err := e.db.FetchTotals(ctx)
		if err != nil {
			return fmt.Errorf("db error: %w", err)
		}
		totals["db"] = t

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error fetching totals: %w", err)
	}

	for k, v := range totals["everef"] {
		if _, ok := totals["db"][k]; !ok {
			log.Printf("missing key %s in db totals", k)
			if err := e.ProcessDay(ctx, k); err != nil {
				return fmt.Errorf("error processing day %s: %w", k, err)
			}
			continue
		}

		if v != totals["db"][k] {
			log.Printf("mismatched value for key %s: %d != %d", k, v, totals["db"][k])
			if err := e.ProcessDay(ctx, k); err != nil {
				return fmt.Errorf("error processing day %s: %w", k, err)
			}
		}
	}

	return nil
}
