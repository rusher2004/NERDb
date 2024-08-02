ALTER TABLE killmail.attacker DROP CONSTRAINT attacker_character_esi_character_id_fk;

ALTER TABLE killmail.attacker DROP CONSTRAINT attacker_esi_killmail_esi_killmail_id_fk;

ALTER TABLE killmail.victim DROP CONSTRAINT victim_character_esi_character_id_fk;

ALTER TABLE killmail.victim DROP CONSTRAINT victim_corporation_esi_corporation_id_fk;
