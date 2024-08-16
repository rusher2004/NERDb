"use server";

import db from "@/app/actions/db";
import { Character } from "@/app/lib/definitions";

export async function getCharacter(id: number): Promise<Character | null> {
  try {
    const [character] = await db<Character[]>`
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
