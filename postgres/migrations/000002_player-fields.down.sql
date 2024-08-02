ALTER TABLE player.alliance DROP COLUMN creator_corporation_id,
  DROP COLUMN creator_id,
  DROP COLUMN date_founded,
  DROP COLUMN executor_corporation_id,
  DROP COLUMN faction_id,
  DROP COLUMN name,
  DROP COLUMN ticker;

ALTER TABLE player.character DROP COLUMN alliance_id,
  DROP COLUMN birthday,
  DROP COLUMN bloodline_id,
  DROP COLUMN corporation_id,
  DROP COLUMN description,
  DROP COLUMN faction,
  DROP COLUMN gender,
  DROP COLUMN name,
  DROP COLUMN race_id,
  DROP COLUMN security_status,
  DROP COLUMN title;

ALTER TABLE player.corporation DROP COLUMN alliance_id,
  DROP COLUMN ceo_id,
  DROP COLUMN creator_id,
  DROP COLUMN date_founded,
  DROP COLUMN description,
  DROP COLUMN faction_id,
  DROP COLUMN home_station_id,
  DROP COLUMN member_count,
  DROP COLUMN name,
  DROP COLUMN shares,
  DROP COLUMN tax_rate,
  DROP COLUMN ticker,
  DROP COLUMN url,
  DROP COLUMN war_eligible;

DROP TYPE player.gender;
