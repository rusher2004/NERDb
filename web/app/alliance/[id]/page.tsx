import Image from "next/image";
import { unstable_cache } from "next/cache";
import { getAlliance } from "@/app/actions/alliance";
import { getCachedAttackersAndVictims } from "@/app/actions/rivals";
import KillmailParticipantCard from "@/app/ui/Cards/KillmailParticipant";

export default async function Page({ params }: { params: { id: string } }) {
  const getCachedAlliance = unstable_cache(
    async (id: number) => await getAlliance(id),
    [`alliance-${params.id}`],
    {
      tags: [`cached-alliance`],
      revalidate: 60 * 60 * 24,
    }
  );
  const alliance = await getCachedAlliance(parseInt(params.id));
  const { attackers, victims } = await getCachedAttackersAndVictims(
    parseInt(params.id),
    "alliance"
  );

  return (
    <div>
      <div className="card card-side card-bordered shadow-xl">
        <figure>
          <Image
            src={`https://images.evetech.net/Alliance/${params.id}_64.png`}
            alt={alliance!.name}
            width={256}
            height={256}
          />
        </figure>
        <div className="card-body">
          <h2 className="card-title">{alliance!.name}</h2>
        </div>
      </div>

      <div className="flex justify-around">
        <div className="join join-vertical">
          <h1>Top Attackers</h1>
          {attackers.map((attacker) => (
            <KillmailParticipantCard
              key={attacker.esiCharacterId}
              {...attacker}
            />
          ))}
        </div>
        <div className="flex flex-col">
          <h1>Top Victims</h1>
          {victims.map((victim) => (
            <KillmailParticipantCard key={victim.esiCharacterId} {...victim} />
          ))}
        </div>
      </div>
    </div>
  );
}
