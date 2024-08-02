ALTER TABLE player.character
  RENAME COLUMN esi_alliance_id TO alliance_id;

ALTER TABLE player.character
  RENAME COLUMN esi_corporation_id TO corporation_id;

ALTER TABLE player.corporation
  RENAME COLUMN esi_alliance_id TO alliance_id;
