ALTER TABLE killmail.attacker
ADD CONSTRAINT attacker_character_character_id_fk FOREIGN KEY (character_id) REFERENCES player.character;

ALTER TABLE killmail.attacker
ADD CONSTRAINT attacker_killmail_killmail_id_fk FOREIGN KEY (killmail_id) REFERENCES killmail.killmail;

ALTER TABLE killmail.victim
ADD CONSTRAINT victim_character_character_id_fk FOREIGN KEY (character_id) REFERENCES player.character;

ALTER TABLE killmail.victim
ADD CONSTRAINT victim_killmail_killmail_id_fk FOREIGN KEY (killmail_id) REFERENCES killmail.killmail;
