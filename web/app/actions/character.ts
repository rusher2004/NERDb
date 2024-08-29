"use server";

import sql from "@/app/lib/db";
import * as Sentry from "@sentry/nextjs";
import { Character } from "@/app/lib/definitions";

export async function getCharacter(id: number): Promise<Character | null> {
  const formData = new FormData();
  formData.append("id", id.toString());

  return Sentry.withServerActionInstrumentation(
    "getCharacter",
    {
      formData: formData,
      recordResponse: true,
    },
    async () => {
      try {
        const [character] = await sql<Character[]>`
          SELECT
            birthday,
            bloodline_id,
            alliance_id,
            character_id,
            corporation_id,
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
            character_id = ${id}
          LIMIT 1
        `;

        return character;
      } catch (err) {
        console.error(err);
        throw new Error(`failed to get character ${id}, ${err}`);
      }
    }
  );
}
