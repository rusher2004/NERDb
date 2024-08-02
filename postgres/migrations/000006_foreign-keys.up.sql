ALTER TABLE killmail.attacker
ADD CONSTRAINT attacker_character_esi_character_id_fk FOREIGN KEY (esi_character_id) REFERENCES player.character;

ALTER TABLE killmail.attacker
ADD CONSTRAINT attacker_esi_killmail_esi_killmail_id_fk FOREIGN KEY (esi_killmail_id) REFERENCES killmail.esi_killmail;

ALTER TABLE killmail.victim
ADD CONSTRAINT victim_character_esi_character_id_fk FOREIGN KEY (esi_character_id) REFERENCES player.character;

ALTER TABLE killmail.victim
ADD CONSTRAINT victim_esi_killmail_esi_killmail_id_fk FOREIGN KEY (esi_killmail_id) REFERENCES killmail.esi_killmail;
