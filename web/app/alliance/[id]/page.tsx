import { getAlliance } from "@/app/actions/alliance";

export default async function Page({ params }: { params: { id: string } }) {
  const alliance = await getAlliance(parseInt(params.id));

  return (
    <div>
      <pre>{JSON.stringify(alliance, null, 2)}</pre>
    </div>
  );
}
