import clsx from "clsx";
import Link from "next/link";
import Avatar from "@/app/ui/Avatar/NewAvatar";
import AllianceInfoBanner from "@/app/ui/Alliance/InfoBanner";
import CorpInfoBanner from "@/app/ui/Corporation/InfoBanner";
import { Character, Participant, ParticipantType } from "@/app/lib/definitions";

function secStatusColor(secStatus: number) {
  return clsx({
    "text-green-500": secStatus > 0,
    "text-red-500": secStatus < 0,
    "text-gray-500": secStatus === 0,
  });
}

export default async function BannerCard({
  subject,
  type,
}: {
  subject: Participant;
  type: ParticipantType;
}) {
  return (
    <div className="card card-side card-bordered shadow-xl">
      <figure>
        <Avatar subject={subject} type={type} />
      </figure>
      <div className="card-body">
        <h2 className="card-title">{subject.name}</h2>

        <div className="flex flex-col gap-1">
          {type === "character" && (subject as Character).securityStatus && (
            <p
              className={clsx("px-2", [
                secStatusColor((subject as Character).securityStatus!),
              ])}
            >
              {(subject as Character)?.securityStatus}
            </p>
          )}
          {type === "character" && (subject as Character).esiCorporationId && (
            <Link
              href={`/corporation/${(subject as Character)!.esiCorporationId}`}
            >
              <CorpInfoBanner id={(subject as Character)!.esiCorporationId} />
            </Link>
          )}
          {type !== "alliance" && subject?.esiAllianceId && (
            <Link href={`/alliance/${subject!.esiAllianceId}`}>
              <AllianceInfoBanner id={subject.esiAllianceId} />
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
