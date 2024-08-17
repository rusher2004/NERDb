import Link from "next/link";
import Image from "next/image";
import { getCorporation } from "@/app/actions/corporation";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import { getCachedAttackersAndVictims } from "@/app/actions/rivals";
import KillmailParticipantCard from "@/app/ui/Cards/KillmailParticipant";

export default async function Page({ params }: { params: { id: string } }) {
  const corporation = await getCorporation(parseInt(params.id));
  const { attackers, victims } = await getCachedAttackersAndVictims(
    parseInt(params.id),
    "corporation"
  );

  return (
    <div>
      <div className="card card-side card-bordered shadow-xl">
        <figure>
          <Image
            src={`https://images.evetech.net/corporations/${corporation?.esiCorporationId}/logo?size=256`}
            alt={corporation!.name}
            width={256}
            height={256}
          />
        </figure>
        <div className="card-body">
          <h2 className="card-title">{corporation!.name}</h2>

          <div className="flex flex-col gap-1">
            {corporation?.esiAllianceId && (
              <Link href={`/alliance/${corporation!.esiAllianceId}`}>
                <AllianceInfoBanner id={corporation.esiAllianceId} />
              </Link>
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
