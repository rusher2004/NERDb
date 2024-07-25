package everef

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func FetchTotals(cl httpClient) (map[string]int, error) {
	req, err := http.NewRequest("GET", "https://data.everef.net/killmails/totals.json", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enexpected status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var totals map[string]int
	if err := json.Unmarshal(b, &totals); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %w", err)
	}

	return totals, nil
}
