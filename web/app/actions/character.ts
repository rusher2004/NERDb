"use server";

import sql from "@/app/actions/db";
import { KillmailParticipant, Character } from "@/app/lib/definitions";

export async function getCharacter(id: number): Promise<Character | null> {
  try {
    const [character] = await sql<Character[]>`
      SELECT
        birthday,
        bloodline_id,
        esi_alliance_id,
        esi_character_id,
        esi_corporation_id,
        description,
        faction_id,
        gender,
        name,
        race_id,
        security_status,
        title
      FROM
        player.character
      WHERE
        esi_character_id = ${id}
      LIMIT 1
    `;

    return character;
  } catch (err) {
    console.error(err);
    return null;
  }
}

export async function getTopAttackers(
  id: number
): Promise<KillmailParticipant[]> {
  try {
    const attackers = await sql<KillmailParticipant[]>`
      WITH killmails AS (
        SELECT DISTINCT
          ka.esi_character_id attacker_id,
          km.esi_killmail_id
        FROM
          player.character vc
          JOIN killmail.victim kv ON kv.esi_character_id = vc.esi_character_id
          JOIN killmail.esi_killmail km ON km.esi_killmail_id = kv.esi_killmail_id
          JOIN killmail.attacker ka ON ka.esi_killmail_id = km.esi_killmail_id
          JOIN player.character ac ON ac.esi_character_id = ka.esi_character_id
        WHERE
          vc.esi_character_id = ${id}
      )
      SELECT
        attacker_id AS esi_character_id,
        COUNT(*)::int AS number_of_kills
      FROM
        killmails
      GROUP BY
        esi_character_id
      ORDER BY
        number_of_kills DESC
      LIMIT 5;
    `;

    return attackers;
  } catch (err) {
    console.error(err);
    return [];
  }
}

export async function getTopAttackersAndVictims(id: number): Promise<{
  attackers: KillmailParticipant[];
  victims: KillmailParticipant[];
}> {
  const [attackers, victims] = await Promise.allSettled([
    getTopAttackers(id),
    getTopVictims(id),
  ]).then((results) => {
    results.forEach((result) => {
      if (result.status === "rejected") {
        console.error(result.reason);
      }
    });

    const out = results.map((result) =>
      result.status === "fulfilled" ? result.value : []
    );

    return out;
  });

  return { attackers, victims };
}

export async function getTopVictims(
  id: number
): Promise<KillmailParticipant[]> {
  try {
    const victims = await sql<KillmailParticipant[]>`
      WITH killmails AS (
        SELECT DISTINCT
          kv.esi_character_id victim_id,
          km.esi_killmail_id
        FROM player.character ac
          JOIN killmail.attacker ka ON ka.esi_character_id = ac.esi_character_id
          JOIN killmail.esi_killmail km ON km.esi_killmail_id = ka.esi_killmail_id
          JOIN killmail.victim kv ON kv.esi_killmail_id = km.esi_killmail_id
          JOIN player.character vc ON vc.esi_character_id = kv.esi_character_id
        WHERE ac.esi_character_id = ${id}
      )
      SELECT 
        victim_id AS esi_character_id,
        COUNT(*)::int AS number_of_kills
      FROM killmails
      GROUP BY esi_character_id
      ORDER BY number_of_kills DESC
      LIMIT 5;
    `;

    return victims;
  } catch (err) {
    console.error(err);
    return [];
  }
}
