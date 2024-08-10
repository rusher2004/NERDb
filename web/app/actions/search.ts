"use server";

import postgres from "postgres";

const db = postgres(process.env.POSTGRES_URL || "", {
  // ssl: process.env.NODE_ENV === "production",
  ssl: false,
});

export interface SearchResult {
  name: string;
  esi_character_id: number;
  esi_corporation_id: number;
  esi_alliance_id?: number;
  rank: number;
}

export async function search(
  query: string
): Promise<{ rows?: SearchResult[]; error?: Error }> {
  if (!query) {
    return { rows: [] };
  }

  const now = new Date().toISOString();
  console.log("searching for", query, "at", now);
  console.time(query + "-" + now);

  try {
    const rows = await db<SearchResult[]>`
    SELECT
      name,
      security_status,
      esi_character_id,
      esi_corporation_id,
      esi_alliance_id,
      ts_rank(to_tsvector('simple', name), websearch_to_tsquery(${query})) rank
    FROM
      player.character
    WHERE to_tsvector('simple', name) @@ websearch_to_tsquery(${query})
    ORDER BY rank desc
    LIMIT 10;
  `;

    console.timeEnd(query + "-" + now);
    console.log(`query ${query}: ${rows.length} results`);

    if (!rows.length) {
      return { rows: [] };
    }

    return { rows: rows };
  } catch (e) {
    console.error(e);
    return { error: new Error(`An error occurred: ${e}`) };
  }
}
