ALTER TABLE killmail.attacker DROP CONSTRAINT attacker_character_character_id_fk;

ALTER TABLE killmail.attacker DROP CONSTRAINT attacker_killmail_killmail_id_fk;

ALTER TABLE killmail.victim DROP CONSTRAINT victim_character_character_id_fk;

ALTER TABLE killmail.victim DROP CONSTRAINT victim_corporation_corporation_id_fk;
