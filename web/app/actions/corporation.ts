"use server";

import sql from "@/app/lib/db";
import { headers } from "next/headers";
import * as Sentry from "@sentry/nextjs";
import { Corporation } from "@/app/lib/definitions";

export async function getCorporation(id: number): Promise<Corporation | null> {
  const formData = new FormData();
  formData.append("id", id.toString());

  return Sentry.withServerActionInstrumentation(
    "getCorporation",
    {
      formData: formData,
      headers: headers(),
      recordResponse: true,
    },
    async () => {
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
        throw new Error(`failed to get corporation ${id}, ${err}`);
      }
    }
  );
}
