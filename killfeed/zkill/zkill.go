package zkill

import (
	"database/sql"
	"encoding/json"

	"github.com/rusher2004/nerdb/killfeed/esi"
)

type JSONNullInt64 struct {
	sql.NullInt64
}

func (i JSONNullInt64) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int64)
	}

	return []byte("null"), nil
}

func (i *JSONNullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.Int64, i.Valid = 0, false
		return nil
	}

	i.Valid = true
	return json.Unmarshal(data, &i.Int64)
}

type RedisQMessage struct {
	Package RedisQPackage `json:"package"`
}

type RedisQPackage struct {
	ESIKill esi.Killmail `json:"killmail"`
	KillID  int          `json:"killID"`
	ZKKill  ZKKillInfo   `json:"zkb"`
}

type ZKKillInfo struct {
	Awox           bool          `json:"awox"`
	DestroyedValue float64       `json:"destroyedValue"`
	DroppedValue   float64       `json:"droppedValue"`
	FittedValue    float64       `json:"fittedValue"`
	Hash           string        `json:"hash"`
	Labels         []string      `json:"labels"`
	LocationID     JSONNullInt64 `json:"locationID"`
	NPC            bool          `json:"npc"`
	Points         int64         `json:"points"`
	Solo           bool          `json:"solo"`
	TotalValue     float64       `json:"totalValue"`
}
