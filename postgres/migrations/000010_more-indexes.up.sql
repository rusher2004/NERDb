-- this one is good.  It's a good one.
CREATE INDEX victim_esi_corporation_id_index ON killmail.victim (esi_corporation_id);

-- note sure on these ones yet.
CREATE INDEX victim_esi_alliance_id_index ON killmail.victim (esi_alliance_id);

CREATE INDEX attacker_esi_killmail_id_index ON killmail.attacker (esi_killmail_id);

CREATE INDEX victim_esi_killmail_id_index ON killmail.victim (esi_killmail_id);
