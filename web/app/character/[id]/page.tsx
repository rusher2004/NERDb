import { getCharacter } from "@/app/actions/character";

export default async function Page({ params }: { params: { id: string } }) {
  const character = await getCharacter(parseInt(params.id));

  return (
    <div>
      <pre>{JSON.stringify(character, null, 2)}</pre>
    </div>
  );
}
