import { getCorporation } from "@/app/actions/corporation";

export default async function Page({ params }: { params: { id: string } }) {
  const corporation = await getCorporation(parseInt(params.id));

  return (
    <div>
      <pre>{JSON.stringify(corporation, null, 2)}</pre>
    </div>
  );
}
