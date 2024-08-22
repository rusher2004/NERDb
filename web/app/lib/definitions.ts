export type Alliance = {
  creatorCorporationID: number;
  creatorId: number;
  dateFounded: Date;
  esiAllianceId: number;
  executorCorporationId?: number;
  factionId?: number;
  name: string;
  ticker: string;
};

export type KillmailParticipant = {
  esiCharacterId: number;
  numberOfKills: number;
};

export type Character = {
  birthday: Date;
  bloodlineId: number;
  esiAllianceId?: number;
  esiCharacterId: number;
  esiCorporationId: number;
  description?: string;
  factionId?: number;
  gender: string;
  name: string;
  raceId: number;
  securityStatus?: number;
  title?: string;
};

export type Corporation = {
  ceoId: number;
  creatorId: number;
  dateFounded?: Date;
  description?: string;
  esiAllianceId?: number;
  esiCorporationId: number;
  factionId?: number;
  homeStationId?: number;
  memberCount: number;
  name: string;
  shares?: number;
  taxRate: number;
  ticker: string;
  url?: string;
  warEligible?: boolean;
};

export type Faction = {
  corporationID: number;
  descriptionID: LocalizedMap;
  iconID: number;
  memberRaces: number[];
  militiaCorporationID: number;
  nameID: LocalizedMap;
  shortDescriptionID: LocalizedMap;
  sizeFactor: number;
  solarSystemID: number;
  uniqueName: boolean;
  factionID: number;
};

export type LocalizedMap = {
  de: string;
  en: string;
  es: string;
  fr: string;
  ja: string;
  ko: string;
  ru: string;
  zh: string;
  it?: string;
};

export type Participant = Alliance | Character | Corporation;

export type ParticipantType =
  | "alliance"
  | "character"
  | "corporation"
  | "faction";
