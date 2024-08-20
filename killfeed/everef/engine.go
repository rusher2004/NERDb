package everef

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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

func NewEngine(hc HTTPClient, cl db.Client) *Engine {
	return &Engine{cl, hc}
}

// openFile attempts to open a file at fp, checking specifically for a file not found error.
func openFile(fp string) (*os.File, error) {
	f, err := os.Open(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %w", err)
		}
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return f, nil
}

// tarFileToSlices is a helper function to copy the contects of f into a slice of T. `in` is
// passed by reference because the function signature is required to just return an error.
func tarFileToSlice[T any](in *[]T) func(ct context.Context, f archiver.File) error {
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

		var km T
		if err := json.Unmarshal(b, &km); err != nil {
			return fmt.Errorf("error unmarshalling killmail: %w", err)
		}

		*in = append(*in, km)

		return nil
	}
}

// ProcessDayKillmails downloads the killmails for a given day from EveRef, then inserts them into
// the database using nerdb/killfeed/db.CopyESIKillmails.
func (e *Engine) ProcessDayKillmails(ctx context.Context, day string) error {
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

	kms := make([]esi.Killmail, 0)
	tar := archiver.Tar{}
	// if err := tar.Extract(ctx, bReader, nil, handleTarFileWithDB(&kms)); err != nil {
	// 	return fmt.Errorf("error extracting tarball: %w", err)
	// }
	if err := tar.Extract(ctx, bReader, nil, tarFileToSlice(&kms)); err != nil {
		return fmt.Errorf("error extracting tarball: %w", err)
	}

	log.Println("extracted kms count:", len(kms))

	if err := e.db.CopyESIKillmails(ctx, day, kms); err != nil {
		return fmt.Errorf("error copying killmails: %w", err)
	}

	return nil
}

// RunCharacterUpdater reads a file of character data and inserts it into the database. A file named
// characters.json is expected to be in the directory provided.
func (e *Engine) RunCharacterUpdater(ctx context.Context, dir string, batchSize int) error {
	// TODO: just use the full file path. Why be cute about looking for it in the directory?
	f, err := openFile(filepath.Join(dir, "characters.json"))
	if err != nil {
		return fmt.Errorf("error opening characters file: %w", err)
	}
	defer f.Close()

	f.Stat()

	scanner := bufio.NewScanner(f)
	maxCapacity := 512 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	chars := make([]db.Character, 0, batchSize)
	batches := 0
	for scanner.Scan() {
		var c Character
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return fmt.Errorf("error unmarshalling character: %w", err)
		}

		chars = append(chars, c.toDBCharacter())

		if (len(chars) == batchSize || !scanner.Scan()) && len(chars) > 0 {
			log.Printf("copying %d characters", batchSize)
			if err := e.db.CopyCharacters(ctx, chars); err != nil {
				return fmt.Errorf("error copying characters: %w", err)
			}

			log.Printf("copied %d total", len(chars)+batchSize*batches)
			batches++

			chars = make([]db.Character, 0, batchSize)
		}

	}

	return nil
}

func (e *Engine) RunCorporationUpdater(ctx context.Context, dir string) error {
	f, err := openFile(filepath.Join(dir, "corporations.json"))
	if err != nil {
		return fmt.Errorf("error opening corporations file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	maxCapacity := 512 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var corps []db.Corporation
	for scanner.Scan() {
		var c Corporation
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return fmt.Errorf("error unmarshalling corporation: %w", err)
		}

		corps = append(corps, c.toDBCorporation())
	}
	log.Printf("read %d corporations", len(corps))

	if err := e.db.CopyCorporations(ctx, corps); err != nil {
		return fmt.Errorf("error copying corporations: %w", err)
	}

	return nil
}

// RunKillmails inserts historical killmail data into the database by comparing the daily totals from
// EveRef and the databse, using ProcessDayKillmails where they differ. Today's and tomorrow's dates
// are skipped to account for timezone differences and allow for Everef to update its dataset.
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

	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	for k, v := range totals["everef"] {
		// do not run today or tomorrow's dates. Depending on what time this is running, EveRef may not
		// yet have the data, and it would throw a 404.
		if k == today.Format("20060102") || k == tomorrow.Format("20060102") {
			log.Printf("skipping day %s", k)
			continue
		}

		if _, ok := totals["db"][k]; !ok {
			log.Printf("missing key %s in db totals", k)
			if err := e.ProcessDayKillmails(ctx, k); err != nil {
				return fmt.Errorf("error processing day %s: %w", k, err)
			}
			continue
		}

		if v != totals["db"][k] {
			log.Printf("mismatched value for key %s: %d != %d", k, v, totals["db"][k])
			if err := e.ProcessDayKillmails(ctx, k); err != nil {
				return fmt.Errorf("error processing day %s: %w", k, err)
			}
		}

		if v == totals["db"][k] {
			log.Printf("day %s is up to date with %d killmails", k, v)
		}
	}

	return nil
}

func (e *Engine) Run(ctx context.Context, dir string, kind ...string) error {
	if !slices.Contains(kind, "character") &&
		!slices.Contains(kind, "corporation") &&
		!slices.Contains(kind, "killmail") {
		return fmt.Errorf("unknown type(s): %s", kind)
	}

	g := new(errgroup.Group)

	if slices.Contains(kind, "character") {
		g.Go(func() error {
			return e.RunCharacterUpdater(ctx, dir, 10_000)
		})
	}

	if slices.Contains(kind, "corporation") {
		g.Go(func() error {
			return e.RunCorporationUpdater(ctx, dir)
		})
	}

	if slices.Contains(kind, "killmail") {
		g.Go(func() error {
			return e.RunKillmails(ctx, e.hc)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error running engine: %w", err)
	}

	return nil
}
