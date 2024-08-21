import { getFaction } from "@/app/actions/faction";
import Avatar from "@/app/ui/Faction/Avatar";

export default async function InfoBanner({ id }: { id: number }) {
  const faction = await getFaction(id);

  return (
    <div className="flex gap-1 justify-start items-start">
      <Avatar id={id} size={64} />
      <div>
        <h1>{faction?.nameID.en}</h1>
      </div>
    </div>
  );
}
