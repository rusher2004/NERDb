ALTER TABLE player.alliance
ADD COLUMN creator_corporation_id INTEGER,
  ADD COLUMN creator_id INTEGER,
  ADD COLUMN date_founded TIMESTAMP,
  ADD COLUMN executor_corporation_id INTEGER,
  ADD COLUMN faction_id INTEGER,
  ADD COLUMN name TEXT,
  ADD COLUMN ticker TEXT;

CREATE TYPE player.gender AS ENUM ('female', 'male');

ALTER TABLE player.character
ADD COLUMN alliance_id INTEGER,
  ADD COLUMN birthday TIMESTAMP,
  ADD COLUMN bloodline_id INTEGER,
  ADD COLUMN corporation_id INTEGER,
  ADD COLUMN description TEXT,
  ADD COLUMN faction_id INTEGER,
  ADD COLUMN gender player.gender,
  ADD COLUMN name TEXT,
  ADD COLUMN race_id INTEGER,
  ADD COLUMN security_status NUMERIC(17, 2),
  ADD COLUMN title TEXT;

ALTER TABLE player.corporation
ADD COLUMN alliance_id INTEGER,
  ADD COLUMN ceo_id INTEGER,
  ADD COLUMN creator_id INTEGER,
  ADD COLUMN date_founded TIMESTAMP,
  ADD COLUMN description TEXT,
  ADD COLUMN faction_id INTEGER,
  ADD COLUMN home_station_id INTEGER,
  ADD COLUMN member_count INTEGER,
  ADD COLUMN name TEXT,
  ADD COLUMN shares INTEGER,
  ADD COLUMN tax_rate NUMERIC(17, 2),
  ADD COLUMN ticker TEXT,
  ADD COLUMN url TEXT,
  ADD COLUMN war_eligible BOOLEAN;
