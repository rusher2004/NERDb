package esi

import "time"

type Killmail struct {
	Attackers     []KillmailAttacker `json:"attackers"`
	KillmailID    int32              `json:"killmail_id"`
	KillmailTime  time.Time          `json:"killmail_time"`
	MoonID        int32              `json:"moon_id,omitempty"`
	SolarSystemID int32              `json:"solar_system_id"`
	Victim        KillmailVictim     `json:"victim"`
	WarID         int32              `json:"war_id,omitempty"`
}

type KillMailParticipant struct {
	AllianceID    int32 `json:"alliance_id,omitempty"`
	CharacterID   int32 `json:"character_id,omitempty"`
	CorporationID int32 `json:"corporation_id,omitempty"`
	FactionID     int32 `json:"faction_id,omitempty"`
	ShipTypeID    int32 `json:"ship_type_id,omitempty"`
}

type KillmailAttacker struct {
	KillMailParticipant

	DamageDone     int32   `json:"damage_done"`
	FactionID      int32   `json:"faction_id,omitempty"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float32 `json:"security_status"`
	WeaponTypeID   int32   `json:"weapon_type_id,omitempty"`
}

type KillmailVictim struct {
	KillMailParticipant

	DamageTaken int32                   `json:"damage_taken"`
	FactionID   int32                   `json:"faction_id,omitempty"`
	Items       []KillmailVictimItem    `json:"items,omitempty"`
	Position    *KillmailVictimPosition `json:"position,omitempty"`
}

type KillmailVictimItem struct {
	Flag              int32                      `json:"flag"`
	ItemTypeID        int32                      `json:"item_type_id"`
	Items             []KillmailVictimItemDetail `json:"items,omitempty"`
	QuantityDestroyed int64                      `json:"quantity_destroyed,omitempty"`
	QuantityDropped   int64                      `json:"quantity_dropped,omitempty"`
	Singleton         int32                      `json:"singleton"`
}

type KillmailVictimItemDetail struct {
	Flag              int32 `json:"flag"`
	ItemTypeID        int32 `json:"item_type_id"`
	QuantityDestroyed int64 `json:"quantity_destroyed,omitempty"`
	QuantityDropped   int64 `json:"quantity_dropped,omitempty"`
	Singleton         int32 `json:"singleton"`
}

type KillmailVictimPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
