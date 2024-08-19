import Link from "next/link";
import { Suspense } from "react";
import Loading from "@/app/ui/Loading/Loading";
import { getCorporation } from "@/app/actions/corporation";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import KillmailParticipants from "@/app/ui/KillmailParticipants/KillmailParticipantList";

export default async function Page({ params }: { params: { id: string } }) {
  const id = parseInt(params.id);
  const corporation = await getCorporation(id);

  return (
    <div className="flex flex-col gap-3">
      <div className="card card-side card-bordered shadow-xl">
        <figure>
          <div className="avatar">
            <div className="max-w-48 max-h-48">
              <img
                src={`https://images.evetech.net/corporations/${corporation?.esiCorporationId}/logo?size=256`}
                alt={corporation!.name}
              />
            </div>
          </div>
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
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type="corporation" role="attacker" />
        </Suspense>
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type="corporation" role="victim" />
        </Suspense>
      </div>
    </div>
  );
}
