"use server";

import sql from "@/app/lib/db";
import * as Sentry from "@sentry/nextjs";
import { unstable_cache } from "next/cache";
import { KillmailParticipant, ParticipantType } from "@/app/lib/definitions";

export const getCachedTopAttackers = unstable_cache(
  async (id: number, type: ParticipantType) => {
    try {
      const formData = new FormData();
      formData.append("id", id.toString());
      formData.append("participantType", type);

      return Sentry.withServerActionInstrumentation(
        "getCachedTopAttackers",
        {
          formData: formData,
          recordResponse: true,
        },
        () => getTopAttackers(id, type)
      );
    } catch (err) {
      console.error(err);
      throw err;
    }
  },
  ["top-attackers"],
  {
    tags: ["cached-top-attackers"],
    revalidate: 60 * 60 * 24,
  }
);

export const getCachedTopVictims = unstable_cache(
  async (id: number, type: ParticipantType) => {
    try {
      const formData = new FormData();
      formData.append("id", id.toString());
      formData.append("participantType", type);

      return Sentry.withServerActionInstrumentation(
        "getCachedTopVictims",
        {
          formData: formData,
          recordResponse: true,
        },
        () => getTopVictims(id, type)
      );
    } catch (err) {
      console.error(err);
      throw err;
    }
  },
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
          return sql`kv.alliance_id = ${id}`;
        case "character":
          return sql`kv.character_id = ${id}`;
        case "corporation":
          return sql`kv.corporation_id = ${id}`;
        case "faction":
          return sql`kv.faction_id = ${id}`;
      }
    };

    const attackers = await sql<KillmailParticipant[]>`
      WITH killmails AS (
        SELECT DISTINCT
          ka.character_id attacker_id,
          km.killmail_id
        FROM
          player.character vc
          JOIN killmail.victim kv ON kv.character_id = vc.character_id
          JOIN killmail.killmail km ON km.killmail_id = kv.killmail_id
          JOIN killmail.attacker ka ON ka.killmail_id = km.killmail_id
          JOIN player.character ac ON ac.character_id = ka.character_id
        WHERE
          ${whereColumn()}
      )
      SELECT
        attacker_id AS character_id,
        COUNT(*)::int AS number_of_kills
      FROM
        killmails
      GROUP BY
        character_id
      ORDER BY
        number_of_kills DESC
      LIMIT 5;
    `;

    return attackers;
  } catch (err) {
    console.error(err);
    throw new Error(`failed to get top attackers for ${type} ${id}, ${err}`);
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
          return sql`ka.alliance_id = ${id}`;
        case "character":
          return sql`ka.character_id = ${id}`;
        case "corporation":
          return sql`ka.corporation_id = ${id}`;
        case "faction":
          return sql`ka.faction_id = ${id}`;
      }
    };

    const victims = await sql<KillmailParticipant[]>`
      WITH killmails AS (
        SELECT DISTINCT
          kv.character_id victim_id,
          km.killmail_id
        FROM player.character ac
          JOIN killmail.attacker ka ON ka.character_id = ac.character_id
          JOIN killmail.killmail km ON km.killmail_id = ka.killmail_id
          JOIN killmail.victim kv ON kv.killmail_id = km.killmail_id
          JOIN player.character vc ON vc.character_id = kv.character_id
        WHERE 
          ${whereColumn()}
      )
      SELECT 
        victim_id AS character_id,
        COUNT(*)::int AS number_of_kills
      FROM killmails
      GROUP BY character_id
      ORDER BY number_of_kills DESC
      LIMIT 5;
    `;

    return victims;
  } catch (err) {
    console.error(err);
    throw new Error(`failed to get top victims for ${type} ${id}, ${err}`);
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
