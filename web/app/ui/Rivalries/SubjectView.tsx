import { Suspense } from "react";
import { unstable_cache } from "next/cache";
import Loading from "@/app/ui/Loading/Loading";
import { getAlliance } from "@/app/actions/alliance";
import BannerCard from "@/app/ui/Rivalries/BannerCard";
import { getCharacter } from "@/app/actions/character";
import { ParticipantType } from "@/app/lib/definitions";
import { getCorporation } from "@/app/actions/corporation";
import KillmailParticipants from "@/app/ui/KillmailParticipants/ParticipantList";

export default async function SubjectView({
  params: { id, type },
}: {
  params: { id: number; type: ParticipantType };
}) {
  const getCachedSubject = unstable_cache(
    async (id: number, type: string) => {
      switch (type) {
        case "character":
          return await getCharacter(id);
        case "corporation":
          return await getCorporation(id);
        case "alliance":
          return await getAlliance(id);
        default:
          return null;
      }
    },
    [`${type}-${id}`],
    {
      tags: [`cached-${type}`],
      revalidate: 60 * 60 * 24,
    }
  );

  const subject = await getCachedSubject(id, type);

  return (
    <div className="flex flex-col gap-3">
      {subject && <BannerCard subject={subject} type={type} />}

      <div className="flex flex-col items-center md:flex-row md:justify-evenly md:items-start gap-6 md:gap-9">
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type={type} role="attacker" />
        </Suspense>
        <Suspense fallback={<Loading />}>
          <KillmailParticipants id={id} type={type} role="victim" />
        </Suspense>
      </div>
    </div>
  );
}
