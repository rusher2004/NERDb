import { Suspense } from "react";
import { unstable_cache } from "next/cache";
import Loading from "@/app/ui/Loading/Loading";
import { getFaction } from "@/app/actions/faction";
import KillmailParticipants from "@/app/ui/KillmailParticipants/ParticipantList";

export default async function Page({ params }: { params: { id: string } }) {
  const id = parseInt(params.id);
  const getCachedFaction = unstable_cache(
    async (id: number) => await getFaction(id),
    [`faction-${params.id}`],
    {
      tags: [`cached-faction`],
      revalidate: 60 * 60 * 24,
    }
  );
  const faction = await getCachedFaction(id);

  return (
    <div className="flex flex-col gap-3">
      <div className="card card-side card-bordered shadow-xl">
        <figure>
          <div className="avatar">
            <div className="max-w-48 max-h-48">
              <img
                src={`https://images.evetech.net/corporations/${id}/logo?size=128`}
                alt={faction!.nameID.en}
              />
            </div>
          </div>
        </figure>
        <div className="card-body">
          <h2 className="card-title">{faction!.nameID.en}</h2>
        </div>
      </div>

      <div className="flex flex-col md:flex-row justify-evenly gap-6 md:gap-9">
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type="faction" role="attacker" />
        </Suspense>
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type="faction" role="victim" />
        </Suspense>
      </div>
    </div>
  );
}
