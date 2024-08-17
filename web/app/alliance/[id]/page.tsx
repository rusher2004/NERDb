import { Suspense } from "react";
import { unstable_cache } from "next/cache";
import Loading from "@/app/ui/Loading/Loading";
import { getAlliance } from "@/app/actions/alliance";
import KillmailParticipants from "@/app/ui/KillmailParticipants/KillmailParticipantList";

export default async function Page({ params }: { params: { id: string } }) {
  const id = parseInt(params.id);
  const getCachedAlliance = unstable_cache(
    async (id: number) => await getAlliance(id),
    [`alliance-${params.id}`],
    {
      tags: [`cached-alliance`],
      revalidate: 60 * 60 * 24,
    }
  );
  const alliance = await getCachedAlliance(id);

  return (
    <div>
      <div className="card card-side card-bordered shadow-xl">
        <figure>
          <div className="avatar">
            <div className="max-w-48 max-h-48">
              <img
                src={`https://images.evetech.net/Alliance/${params.id}_64.png`}
                alt={alliance!.name}
              />
            </div>
          </div>
        </figure>
        <div className="card-body">
          <h2 className="card-title">{alliance!.name}</h2>
        </div>
      </div>

      <div className="flex justify-around">
        <div className="join join-vertical">
          <h1>Top Attackers</h1>
          <Suspense fallback={<Loading />}>
            <KillmailParticipants id={id} type="alliance" side="attacker" />
          </Suspense>
        </div>
        <div className="flex flex-col">
          <h1>Top Victims</h1>
          <Suspense fallback={<Loading />}>
            <KillmailParticipants id={id} type="alliance" side="victim" />
          </Suspense>
        </div>
      </div>
    </div>
  );
}
