import {
  getCorporation,
  getTopAttackersAndVictims,
} from "@/app/actions/corporation";
import Image from "next/image";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import KillmailParticipantCard from "@/app/ui/Cards/KillmailParticipant";

export default async function Page({ params }: { params: { id: string } }) {
  const corporation = await getCorporation(parseInt(params.id));
  const { attackers, victims } = await getTopAttackersAndVictims(
    parseInt(params.id)
  );

  return (
    <div>
      <div className="card lg:card-side card-bordered shadow-xl">
        <figure>
          <Image
            src={`https://images.evetech.net/corporations/${corporation?.esiCorporationId}/logo?size=256`}
            alt={corporation!.name}
            width={256}
            height={256}
            // className="rounded-xl"
          />
        </figure>
        <div className="card-body">
          <h2 className="card-title">{corporation!.name}</h2>

          <div className="flex flex-col gap-1">
            {corporation?.esiAllianceId && (
              <AllianceInfoBanner id={corporation.esiAllianceId} />
            )}
          </div>
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
