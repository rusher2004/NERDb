package esi

import (
	"database/sql"
	"encoding/json"
	"slices"
	"time"
)

type JSONNullInt32 struct {
	sql.NullInt32
}

func (i JSONNullInt32) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int32)
	}

	return []byte("null"), nil
}

func (i *JSONNullInt32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.Int32, i.Valid = 0, false
		return nil
	}

	i.Valid = true
	return json.Unmarshal(data, &i.Int32)
}

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

type Killmail struct {
	Attackers     []KillmailAttacker `json:"attackers"`
	KillmailID    int32              `json:"killmail_id"`
	KillmailTime  time.Time          `json:"killmail_time"`
	MoonID        JSONNullInt32      `json:"moon_id"`
	SolarSystemID int32              `json:"solar_system_id"`
	Victim        KillmailVictim     `json:"victim"`
	WarID         JSONNullInt32      `json:"war_id"`
}

type KillMailParticipant struct {
	AllianceID    JSONNullInt32 `json:"alliance_id"`
	CharacterID   JSONNullInt32 `json:"character_id"`
	CorporationID JSONNullInt32 `json:"corporation_id"`
	FactionID     JSONNullInt32 `json:"faction_id"`
	ShipTypeID    int32         `json:"ship_type_id"`
}

type KillmailAttacker struct {
	KillMailParticipant

	DamageDone     int32         `json:"damage_done"`
	FinalBlow      bool          `json:"final_blow"`
	SecurityStatus float32       `json:"security_status"`
	ShipTypeID     JSONNullInt32 `json:"ship_type_id"`
	WeaponTypeID   JSONNullInt32 `json:"weapon_type_id"`
}

type KillmailVictim struct {
	KillMailParticipant

	DamageTaken int32                   `json:"damage_taken"`
	Items       []KillmailVictimItem    `json:"items,omitempty"`
	Position    *KillmailVictimPosition `json:"position,omitempty"`
}

type KillmailVictimItem struct {
	Flag              int32                 `json:"flag"`
	ItemTypeID        int32                 `json:"item_type_id"`
	Items             *[]KillmailVictimItem `json:"items,omitempty"`
	QuantityDestroyed JSONNullInt64         `json:"quantity_destroyed,omitempty"`
	QuantityDropped   JSONNullInt64         `json:"quantity_dropped,omitempty"`
	Singleton         int32                 `json:"singleton"`
}

type KillmailVictimItemDetail struct {
	Flag              int32         `json:"flag"`
	ItemTypeID        int32         `json:"item_type_id"`
	QuantityDestroyed JSONNullInt64 `json:"quantity_destroyed,omitempty"`
	QuantityDropped   JSONNullInt64 `json:"quantity_dropped,omitempty"`
	Singleton         int32         `json:"singleton"`
}

type KillmailVictimPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (k Killmail) UniqueParticipants() []KillMailParticipant {
	var charIDs []int32
	var participants []KillMailParticipant

	// skip 0 to not include non-unique NPCs
	if k.Victim.KillMailParticipant.CharacterID.Int32 != 0 {
		charIDs = append(charIDs, k.Victim.KillMailParticipant.CharacterID.Int32)
		participants = append(participants, k.Victim.KillMailParticipant)
	}

	for _, a := range k.Attackers {
		if a.CharacterID.Int32 != 0 && !slices.Contains(charIDs, a.CharacterID.Int32) {
			charIDs = append(charIDs, a.CharacterID.Int32)
			participants = append(participants, a.KillMailParticipant)
		}
	}

	return participants
}
