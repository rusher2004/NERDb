import { getCharacter } from "@/app/actions/character";
import { getCachedAttackersAndVictims } from "@/app/actions/rivals";
import KillmailParticipantCard from "@/app/ui/Cards/KillmailParticipant";
import Image from "next/image";
import clsx from "clsx";
import CorpInfoBanner from "@/app/ui/Corporation/InfoBanner";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import Link from "next/link";

function secStatusColor(secStatus: number) {
  return clsx({
    "text-green-500": secStatus > 0,
    "text-red-500": secStatus < 0,
    "text-gray-500": secStatus === 0,
  });
}

export default async function Page({ params }: { params: { id: string } }) {
  const character = await getCharacter(parseInt(params.id));
  const { attackers, victims } = await getCachedAttackersAndVictims(
    parseInt(params.id),
    "character"
  );

  return (
    <div>
      <div className="card card-side card-bordered shadow-xl">
        <figure>
          <Image
            src={`https://images.evetech.net/characters/${character?.esiCharacterId}/portrait?size=256`}
            alt={character!.name}
            width={256}
            height={256}
          />
        </figure>
        <div className="card-body">
          <h2 className="card-title">{character!.name}</h2>

          <div className="flex flex-col gap-1">
            {character?.securityStatus && (
              <p className={secStatusColor(character.securityStatus)}>
                {character?.securityStatus}
              </p>
            )}
            <Link href={`/corporation/${character!.esiCorporationId}`}>
              <CorpInfoBanner id={character!.esiCorporationId} />
            </Link>
            {character?.esiAllianceId && (
              <Link href={`/alliance/${character!.esiAllianceId}`}>
                <AllianceInfoBanner id={character.esiAllianceId} />
              </Link>
            )}
          </div>
          {/* {character?.description && (
            <Markdown className="max-h-32 max-w-32 overflow-scroll">
              {character.description}
            </Markdown>
          )} */}
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
