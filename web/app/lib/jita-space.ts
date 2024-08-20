import { Faction } from "@/app/lib/definitions";

export async function getFaction(id: number): Promise<Faction | null> {
  try {
    const response = await fetch(
      `https://sde.jita.space/latest/universe/factions/${id}`
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch faction ${id}`);
    }

    return await response.json();
  } catch (err) {
    console.error(err);
    return null;
  }
}
