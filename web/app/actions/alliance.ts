"use server";

import sql from "@/app/lib/db";
import * as Sentry from "@sentry/nextjs";
import { Alliance } from "@/app/lib/definitions";

export async function getAlliance(id: number): Promise<Alliance | null> {
  const formData = new FormData();
  formData.append("id", id.toString());

  return Sentry.withServerActionInstrumentation(
    "getAlliance",
    {
      formData: formData,
      recordResponse: true,
    },
    async () => {
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
          alliance_id = ${id}
        LIMIT 1
      `;

        return alliance;
      } catch (err) {
        console.error(err);
        throw new Error(`failed to get alliance ${id}, ${err}`);
      }
    }
  );
}
