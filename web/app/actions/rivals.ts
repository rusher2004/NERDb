"use server";

import sql from "@/app/lib/db";
import { unstable_cache } from "next/cache";
import { KillmailParticipant, ParticipantType } from "@/app/lib/definitions";

export const getCachedTopAttackers = unstable_cache(
  async (id: number, type: ParticipantType) => await getTopAttackers(id, type),
  ["top-attackers"],
  {
    tags: ["cached-top-attackers"],
    revalidate: 60 * 60 * 24,
  }
);

export const getCachedTopVictims = unstable_cache(
  async (id: number, type: ParticipantType) => await getTopVictims(id, type),
  ["top-victims"],
  {
    tags: ["cached-top-victims"],
    revalidate: 60 * 60 * 24,
  }
);

export const getCachedAttackersAndVictims = unstable_cache(
  async (id: number, type: ParticipantType) =>
    await getTopAttackersAndVictims(id, type),
  ["attackers-and-victims"],
  {
    tags: ["cached-attackers-and-victims"],
    revalidate: 60 * 60 * 24,
  }
);

async function getTopAttackers(
  id: number,
  type: ParticipantType
): Promise<KillmailParticipant[]> {
  try {
    const whereColumn = () => {
      switch (type) {
        case "alliance":
          return sql`vc.esi_alliance_id = ${id}`;
        case "character":
          return sql`vc.esi_character_id = ${id}`;
        case "corporation":
          return sql`vc.esi_corporation_id = ${id}`;
      }
    };

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
          ${whereColumn()}
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

async function getTopVictims(
  id: number,
  type: ParticipantType
): Promise<KillmailParticipant[]> {
  try {
    const whereColumn = () => {
      switch (type) {
        case "alliance":
          return sql`ac.esi_alliance_id = ${id}`;
        case "character":
          return sql`ac.esi_character_id = ${id}`;
        case "corporation":
          return sql`ac.esi_corporation_id = ${id}`;
      }
    };

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
        WHERE 
          ${whereColumn()}
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

async function getTopAttackersAndVictims(
  id: number,
  type: ParticipantType
): Promise<{
  attackers: KillmailParticipant[];
  victims: KillmailParticipant[];
}> {
  const [attackers, victims] = await Promise.allSettled([
    getTopAttackers(id, type),
    getTopVictims(id, type),
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
