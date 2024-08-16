CREATE INDEX attacker_esi_killmail_id_esi_character_id_index ON killmail.attacker (esi_killmail_id, esi_character_id);

CREATE INDEX victim_esi_killmail_id_esi_character_id_index ON killmail.victim (esi_killmail_id, esi_character_id);

CREATE INDEX attacker_esi_character_id_index ON killmail.attacker (esi_character_id);

CREATE INDEX victim_esi_character_id_index ON killmail.victim (esi_character_id);
