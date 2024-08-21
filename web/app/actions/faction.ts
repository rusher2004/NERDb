"use server";

import sql from "@/app/lib/db";
import { Faction } from "@/app/lib/definitions";
import { getFaction as fetchFaction } from "@/app/lib/jita-space";

export async function getFaction(id: number): Promise<Faction | null> {
  return await fetchFaction(id);
}
