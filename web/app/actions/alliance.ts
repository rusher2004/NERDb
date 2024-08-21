"use server";

import sql from "@/app/lib/db";
import { Alliance } from "@/app/lib/definitions";

export async function getAlliance(id: number): Promise<Alliance | null> {
  try {
    const [alliance] = await sql<Alliance[]>`
      SELECT
        creator_corporation_id,
        creator_id,
        date_founded,
        executor_corporation_id,
        faction_id,
        name,
        ticker
      FROM
        player.alliance
      WHERE
        esi_alliance_id = ${id}
      LIMIT 1
    `;

    return alliance;
  } catch (err) {
    console.error(err);
    return null;
  }
}
