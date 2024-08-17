import { getCorporation } from "@/app/actions/corporation";
import Avatar from "@/app/ui/Corporation/Avatar";

export default async function InfoBanner({ id }: { id: number }) {
  const corporation = await getCorporation(id);

  return (
    <div className="flex gap-1">
      <Avatar id={id} size={64} />
      <div>
        <h1>{corporation?.name}</h1>
        <p>{corporation?.ticker}</p>
      </div>
    </div>
  );
}
