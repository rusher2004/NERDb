package db

import "github.com/rusher2004/nerdb/esi"

func FetchTotals() (map[string]int, error) {
	return map[string]int{"20121113": 11983}, nil
}

func UpsertKillmail(km esi.Killmail) error {
	return nil
}
