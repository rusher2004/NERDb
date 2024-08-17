import clsx from "clsx";
import Link from "next/link";
import Image from "next/image";
import { Suspense } from "react";
import { getCharacter } from "@/app/actions/character";
import CorpInfoBanner from "@/app/ui/Corporation/InfoBanner";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import KillmailParticipants from "@/app/ui/KillmailParticipants/KillmailParticipantList";

function secStatusColor(secStatus: number) {
  return clsx({
    "text-green-500": secStatus > 0,
    "text-red-500": secStatus < 0,
    "text-gray-500": secStatus === 0,
  });
}

export default async function Page({ params }: { params: { id: string } }) {
  const id = parseInt(params.id);
  const character = await getCharacter(id);

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
          <Suspense fallback={<div>Loading...</div>}>
            <KillmailParticipants id={id} type="character" side="attacker" />
          </Suspense>
        </div>
        <div className="flex flex-col">
          <h1>Top Victims</h1>
          <Suspense fallback={<div>Loading...</div>}>
            <KillmailParticipants id={id} type="character" side="victim" />
          </Suspense>
        </div>
      </div>
    </div>
  );
}
