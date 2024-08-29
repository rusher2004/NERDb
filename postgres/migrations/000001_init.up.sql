/*
 functions
 */
CREATE OR REPLACE FUNCTION updated_timestamp() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at := CURRENT_TIMESTAMP;

RETURN NEW;

END;

$$ LANGUAGE PLPGSQL;

/*
 player schema
 */
CREATE SCHEMA IF NOT EXISTS player;

CREATE TABLE IF NOT EXISTS player.character (
  character_id integer NOT NULL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER player_character_updated_trigger BEFORE
UPDATE ON player.character FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS player.corporation (
  corporation_id integer NOT NULL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER player_corporation_updated_trigger BEFORE
UPDATE ON player.corporation FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS player.alliance (
  alliance_id integer NOT NULL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER player_alliance_updated_trigger BEFORE
UPDATE ON player.alliance FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

/*
 killmail schema
 */
CREATE SCHEMA IF NOT EXISTS killmail;

/*
 types
 */
DO $$ BEGIN CREATE TYPE killmail.victim_item_parent_type AS enum ('victim', 'item');

EXCEPTION
WHEN duplicate_object THEN NULL;

END $$;

CREATE TABLE IF NOT EXISTS killmail.killmail (
  killmail_id integer NOT NULL CONSTRAINT killmail_pk PRIMARY KEY,
  time timestamp NOT NULL,
  moon_id integer,
  solar_system_id integer NOT NULL,
  war_id integer,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER killmail_killmail_updated_trigger BEFORE
UPDATE ON killmail.killmail FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS killmail.zkill_info (
  zkill_info_id serial PRIMARY KEY,
  killmail_id integer NOT NULL CONSTRAINT zkill_info_pk UNIQUE REFERENCES killmail.killmail(killmail_id),
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
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER killmail_zkill_info_updated_trigger BEFORE
UPDATE ON killmail.zkill_info FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS killmail.victim (
  character_id integer REFERENCES player.character(character_id),
  killmail_id integer NOT NULL REFERENCES killmail.killmail(killmail_id),
  damage_taken integer NOT NULL,
  alliance_id integer,
  corporation_id integer,
  faction_id integer,
  position_x double precision,
  position_y double precision,
  position_z double precision,
  ship_type_id integer NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER killmail_victim_updated_trigger BEFORE
UPDATE ON killmail.victim FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS killmail.attacker (
  character_id integer REFERENCES player.character(character_id),
  killmail_id integer NOT NULL REFERENCES killmail.killmail(killmail_id),
  damage_done integer NOT NULL,
  alliance_id integer,
  corporation_id integer,
  faction_id integer,
  final_blow boolean NOT NULL,
  security_status double precision,
  ship_type_id integer,
  weapon_type_id integer,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER killmail_attacker_updated_trigger BEFORE
UPDATE ON killmail.attacker FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

CREATE TABLE IF NOT EXISTS killmail.victim_item (
  victim_item_id serial PRIMARY KEY,
  parent_id integer NOT NULL,
  parent_type killmail.victim_item_parent_type NOT NULL,
  flag integer NOT NULL,
  item_type_id integer NOT NULL,
  quantity_destroyed integer,
  quantity_dropped integer,
  singleton integer NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER killmail_victim_item_updated_trigger BEFORE
UPDATE ON killmail.victim_item FOR EACH ROW EXECUTE FUNCTION updated_timestamp();

/*
 universe schema
 */
CREATE SCHEMA IF NOT EXISTS universe;

CREATE TABLE IF NOT EXISTS universe.faction (
  faction_id integer PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER universe_faction_updated_trigger BEFORE
UPDATE ON universe.faction FOR EACH ROW EXECUTE FUNCTION updated_timestamp();
