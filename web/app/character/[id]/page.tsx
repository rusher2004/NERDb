import clsx from "clsx";
import Link from "next/link";
import { Suspense } from "react";
import Loading from "@/app/ui/Loading/Loading";
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
    <div className="flex flex-col gap-3">
      <div className="card card-side card-compact card-bordered shadow-xl">
        <figure>
          <div className="avatar">
            <div className="max-w-48 max-h-48">
              <img
                src={`https://images.evetech.net/characters/${character?.esiCharacterId}/portrait?size=256`}
                alt={character!.name}
              />
            </div>
          </div>
        </figure>
        <div className="card-body">
          <h2 className="card-title">{character!.name}</h2>

          {/* {character?.securityStatus && (
            <p className={secStatusColor(character.securityStatus)}>
              {character?.securityStatus}
            </p>
          )} */}

          <div className="flex flex-col gap-1">
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
      <div className="flex justify-evenly gap-6 md:gap-9">
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type="character" role="attacker" />
        </Suspense>
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type="character" role="victim" />
        </Suspense>
      </div>
    </div>
  );
}
