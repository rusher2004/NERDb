"use server";

import sql from "@/app/lib/db";
import { Corporation } from "@/app/lib/definitions";

export async function getCorporation(id: number): Promise<Corporation | null> {
  try {
    const [corporation] = await sql<Corporation[]>`
      SELECT
        ceo_id,
        creator_id,
        date_founded,
        description,
        esi_alliance_id,
        esi_corporation_id,
        faction_id,
        home_station_id,
        member_count,
        name,
        shares,
        tax_rate,
        ticker,
        url,
        war_eligible
      FROM
        player.corporation
      WHERE
        esi_corporation_id = ${id}
      LIMIT 1
    `;

    return corporation;
  } catch (err) {
    console.error(err);
    return null;
  }
}
