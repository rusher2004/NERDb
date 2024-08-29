ALTER TABLE player.character
  RENAME COLUMN alliance_id TO alliance_id;

ALTER TABLE player.character
  RENAME corporation_id TO corporation_id;

ALTER TABLE player.corporation
  RENAME COLUMN alliance_id TO alliance_id;
