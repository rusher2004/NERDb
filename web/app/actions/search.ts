"use server";

import sql from "@/app/lib/db";
import { headers } from "next/headers";
import * as Sentry from "@sentry/nextjs";

export interface SearchAllResult {
  name: string;
  id: number;
  esiCorporationId?: number;
  esiAllianceId?: number;
  rank: number;
  type: "character" | "corporation" | "alliance";
}

export async function searchAll(
  query: string
): Promise<{ rows?: SearchAllResult[]; error?: Error }> {
  if (!query) {
    return { rows: [] };
  }

  const formData = new FormData();
  formData.append("query", query);

  const now = new Date().toISOString();
  console.log("searching for all", query, "at", now);
  console.time(query + "-" + now);

  return Sentry.withServerActionInstrumentation(
    "searchAll",
    {
      formData: formData,
      headers: headers(),
      recordResponse: true,
    },
    async () => {
      try {
        const rows = await sql<SearchAllResult[]>`
    WITH matched_chars AS (
      SELECT
        name,
        esi_character_id AS id,
        esi_corporation_id,
        esi_alliance_id,
        'character' AS type,
        ts_rank(to_tsvector('simple',
            name),
          websearch_to_tsquery (${query})) AS rank
      FROM
        player.character
      WHERE
        to_tsvector('simple',
          name) @@ websearch_to_tsquery (${query})
      ORDER BY
        rank DESC,
        name
      LIMIT 15
    ),
    matched_corps AS (
      SELECT
        name,
        esi_corporation_id AS id,
        esi_corporation_id,
        esi_alliance_id,
        'corporation' AS type,
        ts_rank(to_tsvector('simple',
            name),
          websearch_to_tsquery (${query})) AS rank
      FROM
        player.corporation
      WHERE
        to_tsvector('simple',
          name) @@ websearch_to_tsquery (${query})
      ORDER BY
        rank DESC,
        name
      LIMIT 15
    ),
    matched_allis AS (
      SELECT
        name,
        esi_alliance_id AS id,
        -1 as esi_corporation_id,
        esi_alliance_id,
        'alliance' AS type,
        ts_rank(to_tsvector('simple',
            name),
          websearch_to_tsquery (${query})) AS rank
      FROM
        player.alliance
      WHERE
        to_tsvector('simple',
          name) @@ websearch_to_tsquery (${query})
      ORDER BY
        rank DESC,
        name
      LIMIT 15
    )
    SELECT name, id, esi_corporation_id, esi_alliance_id, rank, type
    FROM
      matched_chars
    UNION
    SELECT name, id, esi_corporation_id, esi_alliance_id, rank, type
    FROM
      matched_corps
    UNION
    SELECT name, id, esi_corporation_id, esi_alliance_id, rank, type
    FROM
      matched_allis
    ORDER BY
      rank DESC,
      name
    LIMIT 15;
    `;

        console.timeEnd(query + "-" + now);
        console.log(`query ${query}: ${rows.length} results`);

        if (!rows.length) {
          return { rows: [] };
        }

        return { rows: rows };
      } catch (err) {
        console.error(err);
        throw new Error(`search failed: ${err}`);
      }
    }
  );
}
