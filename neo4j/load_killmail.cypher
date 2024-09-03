LOAD CSV WITH HEADERS FROM 'file:///killmail_killmail.csv.gz' AS row
CALL (row) {
MERGE (k:Killmail {killmailId: toInteger(row.killmail_id)})
ON CREATE SET
  // 2024-08-29 23:10:48.146386+00
  k.time = datetime({epochMillis: apoc.date.parse(row.time, "ms", "yyyy-MM-dd HH:mm:ssX")}),
  k.moonId = toInteger(row.moon_id),
  k.solarSystemId = toInteger(row.solar_system_id),
  k.warId = toInteger(row.war_id),
  k.created = datetime()
MERGE (v:Victim)-[:VICTIM_OF]->(k)
ON CREATE SET
  v.characterId = toInteger(nullIf(trim(row.character_id), "")),
  v.corporationId = toInteger(nullIf(trim(row.corporation_id), "")),
  v.allianceId = toInteger(nullIf(trim(row.alliance_id), "")),
  v.factionId = toInteger(nullIf(trim(row.faction_id), "")),
  v.damageTaken = toInteger(row.damage_taken),
  v.position = point({x: toFloat(row.position_x), y: toFloat(row.position_y), z: toFloat(row.position_z)}),
  v.shipTypeId = toInteger(row.ship_type_id),
  v.created = datetime()
} IN TRANSACTIONS OF 250000 ROWS;

LOAD CSV WITH HEADERS FROM 'file:///killmail_attacker.csv.gz' AS row
CALL (row) {
MATCH (k:Killmail {killmailId: toInteger(row.killmail_id)})
WITH k, apoc.map.clean({
  characterId: toInteger(nullIf(trim(row.character_id), "")),
  corporationId: toInteger(nullIf(trim(row.corporation_id), "")),
  allianceId: toInteger(nullIf(trim(row.alliance_id), "")),
  factionId: toInteger(nullIf(trim(row.faction_id), "")),
  damageDone: toInteger(row.damage_done),
  finalBlow: CASE WHEN row.final_blow = 't' THEN true ELSE NULL END,
  securityStatus: toFloat(row.security_status),
  shipTypeId: toInteger(row.ship_type_id),
  weaponTypeId: toInteger(row.weapon_type_id),
  created: datetime()
}, [], [null]) as cleaned
CALL apoc.create.node(
    ["Attacker"],
    cleaned
) YIELD node
CREATE (node)-[:ATTACKER_OF]->(k)
} IN TRANSACTIONS OF 100000 ROWS;
