/*
 functions
 */
CREATE OR REPLACE FUNCTION updated_timestamp() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at := CURRENT_TIMESTAMP;

RETURN NEW;

END;

$$ LANGUAGE PLPGSQL;

/*
 universe schema
 */
CREATE SCHEMA IF NOT EXISTS universe;

CREATE TABLE IF NOT EXISTS universe.faction (
  faction_id integer NOT NULL PRIMARY KEY,
  corporation_id INTEGER,
  description TEXT NOT NULL,
  is_unique BOOLEAN NOT NULL,
  militia_corporation_id INTEGER,
  name TEXT NOT NULL,
  size_factor NUMERIC(17, 2) NOT NULL,
  solar_system_id INTEGER,
  station_count INTEGER NOT NULL,
  station_system_count INTEGER NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

/*
 player schema
 */
CREATE SCHEMA IF NOT EXISTS player;

CREATE TYPE player.gender AS ENUM ('female', 'male', '');

CREATE TABLE IF NOT EXISTS player.character (
  character_id integer NOT NULL PRIMARY KEY,
  alliance_id INTEGER,
  birthday TIMESTAMPTZ,
  bloodline_id INTEGER,
  corporation_id INTEGER,
  description TEXT,
  deleted boolean,
  faction_id INTEGER,
  gender player.gender,
  name TEXT,
  race_id INTEGER,
  security_status NUMERIC(17, 2),
  title TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER player_character_updated_trigger BEFORE
UPDATE ON player.character FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS player.corporation (
  corporation_id integer NOT NULL PRIMARY KEY,
  alliance_id INTEGER,
  ceo_id INTEGER,
  creator_id INTEGER,
  date_founded TIMESTAMPTZ,
  description TEXT,
  faction_id INTEGER,
  home_station_id INTEGER,
  member_count INTEGER,
  name TEXT,
  shares BIGINT,
  tax_rate NUMERIC(17, 2),
  ticker TEXT,
  url TEXT,
  war_eligible BOOLEAN,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER player_corporation_updated_trigger BEFORE
UPDATE ON player.corporation FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS player.alliance (
  alliance_id integer NOT NULL PRIMARY KEY,
  creator_corporation_id INTEGER,
  creator_id INTEGER,
  date_founded TIMESTAMPTZ,
  executor_corporation_id INTEGER,
  faction_id INTEGER,
  name TEXT,
  ticker TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER player_alliance_updated_trigger BEFORE
UPDATE ON player.alliance FOR EACH ROW EXECUTE FUNCTION updated_timestamp();


/*
 foreign keys
 */
-- we can only guarantee that the faction is in the database
ALTER TABLE player.character ADD CONSTRAINT character_faction_fk FOREIGN KEY (faction_id) REFERENCES universe.faction(faction_id);
-- ALTER TABLE player.character ADD CONSTRAINT character_alliance_fk FOREIGN KEY (alliance_id) REFERENCES player.alliance(alliance_id);
-- ALTER TABLE player.character ADD CONSTRAINT character_corporation_fk FOREIGN KEY (corporation_id) REFERENCES player.corporation(corporation_id);

ALTER TABLE player.corporation ADD CONSTRAINT corporation_faction_fk FOREIGN KEY (faction_id) REFERENCES universe.faction(faction_id);
-- ALTER TABLE player.corporation ADD CONSTRAINT corporation_alliance_fk FOREIGN KEY (alliance_id) REFERENCES player.alliance(alliance_id);
-- ALTER TABLE player.corporation ADD CONSTRAINT corporation_ceo_fk FOREIGN KEY (ceo_id) REFERENCES player.character(character_id);
-- ALTER TABLE player.corporation ADD CONSTRAINT corporation_creator_fk FOREIGN KEY (creator_id) REFERENCES player.character(character_id);

ALTER TABLE player.alliance ADD CONSTRAINT alliance_faction_fk FOREIGN KEY (faction_id) REFERENCES universe.faction(faction_id);
-- ALTER TABLE player.alliance ADD CONSTRAINT alliance_creator_corporation_fk FOREIGN KEY (creator_corporation_id) REFERENCES player.corporation(corporation_id);
-- ALTER TABLE player.alliance ADD CONSTRAINT alliance_creator_fk FOREIGN KEY (creator_id) REFERENCES player.character(character_id);
-- ALTER TABLE player.alliance ADD CONSTRAINT alliance_executor_corporation_fk FOREIGN KEY (executor_corporation_id) REFERENCES player.corporation(corporation_id);

/*
 killmail schema
 */
CREATE SCHEMA IF NOT EXISTS killmail;

CREATE TABLE IF NOT EXISTS killmail.killmail (
  killmail_id integer NOT NULL,
  time TIMESTAMPTZ NOT NULL,
  moon_id integer,
  solar_system_id integer NOT NULL,
  war_id integer,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

  -- victim info
  character_id integer,
  alliance_id integer,
  corporation_id integer,
  faction_id integer,
  damage_taken integer,
  position_x double precision,
  position_y double precision,
  position_z double precision,
  ship_type_id integer
);

SELECT create_hypertable('killmail.killmail', by_range('time'));

/*
 Timescale Indexes
 */
CREATE INDEX killmail_killmail_id_time_idx ON killmail.killmail (killmail_id, time DESC);
CREATE INDEX killmail_alliance_id_time_idx ON killmail.killmail (alliance_id, time DESC)
  WHERE alliance_id IS NOT NULL;
CREATE INDEX killmail_character_id_time_idx ON killmail.killmail (character_id, time DESC)
  WHERE character_id IS NOT NULL;
CREATE INDEX killmail_corporation_id_time_idx ON killmail.killmail (corporation_id, time DESC)
  WHERE corporation_id IS NOT NULL;
CREATE INDEX killmail_faction_id_time_idx ON killmail.killmail (faction_id, time DESC)
  WHERE faction_id IS NOT NULL;

CREATE OR REPLACE TRIGGER killmail_killmail_updated_trigger BEFORE
UPDATE ON killmail.killmail FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS killmail.zkill_info (
  zkill_info_id serial PRIMARY KEY,
  killmail_id integer NOT NULL,
  awox boolean NOT NULL,
  destroyed_value numeric(17, 2) NOT NULL,
  dropped_value numeric(17, 2) NOT NULL,
  fitted_value numeric(17, 2) NOT NULL,
  hash text NOT NULL,
  location_id integer,
  npc boolean NOT NULL,
  points integer NOT NULL,
  solo boolean NOT NULL,
  total_value numeric(17, 2) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX zkill_info_killmail_id_idx ON killmail.zkill_info (killmail_id);

CREATE OR REPLACE TRIGGER killmail_zkill_info_updated_trigger BEFORE
UPDATE ON killmail.zkill_info FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS killmail.attacker (
  character_id integer,
  killmail_id integer NOT NULL,
  damage_done integer NOT NULL,
  alliance_id integer,
  corporation_id integer,
  faction_id integer,
  final_blow boolean NOT NULL,
  security_status double precision,
  ship_type_id integer,
  weapon_type_id integer,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX attacker_killmail_id_idx ON killmail.attacker (killmail_id);
CREATE INDEX attacker_alliance_id_idx ON killmail.attacker (alliance_id);
CREATE INDEX attacker_character_id_idx ON killmail.attacker (character_id);
CREATE INDEX attacker_corporation_id_idx ON killmail.attacker (corporation_id);
CREATE INDEX attacker_faction_id_idx ON killmail.attacker (faction_id);

CREATE OR REPLACE TRIGGER killmail_attacker_updated_trigger BEFORE
UPDATE ON killmail.attacker FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

/*
 Search indexes
 */

CREATE INDEX character_name_search_idx ON player.character USING GIN (to_tsvector('simple', name));

CREATE INDEX corporation_name_search_idx ON player.corporation USING GIN (to_tsvector('simple', name));

CREATE INDEX alliance_name_search_idx ON player.alliance USING GIN (to_tsvector('simple', name));
