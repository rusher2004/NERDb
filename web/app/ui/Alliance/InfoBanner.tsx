import { getAlliance } from "@/app/actions/alliance";
import Avatar from "@/app/ui/Alliance/Avatar";

export default async function InfoBanner({ id }: { id: number }) {
  const alliance = await getAlliance(id);

  return (
    <div className="flex gap-1 justify-start items-center">
      <Avatar id={id} size={64} />
      <div>
        <h1>{alliance?.name}</h1>
        <p>{alliance?.ticker}</p>
      </div>
    </div>
  );
}
