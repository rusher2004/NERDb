export type Alliance = {
  creatorCorporationID: number;
  creatorID: number;
  dateFounded: Date;
  executorCorporationID?: number;
  factionID?: number;
  name: string;
  ticker: string;
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
