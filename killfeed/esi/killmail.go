package esi

import (
	"slices"
	"time"

	"github.com/rusher2004/nerdb/null"
)

type Killmail struct {
	Attackers     []KillmailAttacker `json:"attackers"`
	KillmailID    int32              `json:"killmail_id"`
	KillmailTime  time.Time          `json:"killmail_time"`
	MoonID        null.JSONNullInt32 `json:"moon_id"`
	SolarSystemID int32              `json:"solar_system_id"`
	Victim        KillmailVictim     `json:"victim"`
	WarID         null.JSONNullInt32 `json:"war_id"`
}

type KillMailParticipant struct {
	AllianceID    null.JSONNullInt32 `json:"alliance_id"`
	CharacterID   null.JSONNullInt32 `json:"character_id"`
	CorporationID null.JSONNullInt32 `json:"corporation_id"`
	FactionID     null.JSONNullInt32 `json:"faction_id"`
	ShipTypeID    int32              `json:"ship_type_id"`
}

type KillmailAttacker struct {
	KillMailParticipant

	DamageDone     int32              `json:"damage_done"`
	FinalBlow      bool               `json:"final_blow"`
	SecurityStatus float32            `json:"security_status"`
	ShipTypeID     null.JSONNullInt32 `json:"ship_type_id"`
	WeaponTypeID   null.JSONNullInt32 `json:"weapon_type_id"`
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
	QuantityDestroyed null.JSONNullInt64    `json:"quantity_destroyed,omitempty"`
	QuantityDropped   null.JSONNullInt64    `json:"quantity_dropped,omitempty"`
	Singleton         int32                 `json:"singleton"`
}

type KillmailVictimItemDetail struct {
	Flag              int32              `json:"flag"`
	ItemTypeID        int32              `json:"item_type_id"`
	QuantityDestroyed null.JSONNullInt64 `json:"quantity_destroyed,omitempty"`
	QuantityDropped   null.JSONNullInt64 `json:"quantity_dropped,omitempty"`
	Singleton         int32              `json:"singleton"`
}

type KillmailVictimPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// UniqueParticipants returns a list of unique participants in the killmail, combined from
// attackers and victim.
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
