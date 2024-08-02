ALTER TABLE player.character
  RENAME COLUMN alliance_id TO esi_alliance_id;

ALTER TABLE player.character
  RENAME corporation_id TO esi_corporation_id;

ALTER TABLE player.corporation
  RENAME COLUMN alliance_id TO esi_alliance_id;
