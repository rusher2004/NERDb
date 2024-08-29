CREATE INDEX attacker_killmail_id_character_id_index ON killmail.attacker (killmail_id, character_id);

CREATE INDEX victim_killmail_id_character_id_index ON killmail.victim (killmail_id, character_id);

CREATE INDEX attacker_character_id_index ON killmail.attacker (character_id);

CREATE INDEX victim_character_id_index ON killmail.victim (character_id);

CREATE INDEX attacker_faction_id_index ON killmail.attacker (faction_id);

CREATE INDEX victim_faction_id_index ON killmail.victim (faction_id);
