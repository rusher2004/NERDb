import Link from "next/link";
import Image from "next/image";
import { Suspense } from "react";
import { getCorporation } from "@/app/actions/corporation";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import KillmailParticipants from "@/app/ui/KillmailParticipants/KillmailParticipantList";

export default async function Page({ params }: { params: { id: string } }) {
  const id = parseInt(params.id);
  const corporation = await getCorporation(id);

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
          <Suspense fallback={<div>Loading...</div>}>
            <KillmailParticipants id={id} type="corporation" side="attacker" />
          </Suspense>
        </div>
        <div className="flex flex-col">
          <h1>Top Victims</h1>
          <Suspense fallback={<div>Loading...</div>}>
            <KillmailParticipants id={id} type="corporation" side="victim" />
          </Suspense>
        </div>
      </div>
    </div>
  );
}
