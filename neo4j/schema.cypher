CREATE CONSTRAINT killmail_id_unique
FOR (killmail:Killmail) REQUIRE killmail.killmailId IS UNIQUE;

// CREATE CONSTRAINT killmail_id_required
// FOR (killmail:Killmail) REQUIRE killmail.killmailId IS NOT null;

CREATE INDEX victim_character_id
FOR (victim:Victim) ON (victim.characterId);

CREATE INDEX victim_corporation_id
FOR (victim:Victim) ON (victim.corporationId);

CREATE INDEX victim_alliance_id
FOR (victim:Victim) ON (victim.allianceId);

CREATE INDEX victim_faction_id
FOR (victim:Victim) ON (victim.factionId);

CREATE INDEX attacker_character_id
FOR (attacker:Attacker) ON (attacker.characterId);

CREATE INDEX attacker_corporation_id
FOR (attacker:Attacker) ON (attacker.corporationId);

CREATE INDEX attacker_alliance_id
FOR (attacker:Attacker) ON (attacker.allianceId);

CREATE INDEX attacker_faction_id
FOR (attacker:Attacker) ON (attacker.factionId);
