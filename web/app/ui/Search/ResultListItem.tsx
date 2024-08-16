import Avatar from "@/app/ui/Avatar/Avatar";
import AllianceAvatar from "@/app/ui/Alliance/Avatar";
import CharAvatar from "@/app/ui/Character/Avatar";
import CorpAvatar from "@/app/ui/Corporation/Avatar";
import { clsx } from "clsx";
import Link from "next/link";

export interface Props {
  id: number;
  name: string;
  esiCorporationId?: number;
  esiAllianceId?: number;
  type: "character" | "corporation" | "alliance";
}

function resultAvatar(result: Props) {
  if (result.type === "character") {
    return <CharAvatar size={32} id={result.id} />;
  }

  if (result.type === "corporation") {
    return <CorpAvatar size={32} id={result.id} />;
  }

  if (result.type === "alliance") {
    return <AllianceAvatar size={32} id={result.id} />;
  }

  return null;
}

function secStatusColor(secStatus: number) {
  return clsx({
    "text-green-500": secStatus > 0,
    "text-red-500": secStatus < 0,
    "text-gray-500": secStatus === 0,
  });
}

export default function ResultListItem(props: Props) {
  console.log("ResultListItem", props);
  return (
    <Link
      href={`/${props.type}/${props.id}`}
      prefetch={false}
      className="flex justify-between w-full"
    >
      <div className="flex gap-1 items-start">
        {resultAvatar(props)}
        <div className="flex flex-col gap-1">
          <h2>{props.name}</h2>
        </div>
      </div>
      <div className="flex gap-1">
        {props.esiCorporationId && props.type === "character" && (
          <CorpAvatar size={32} id={props.esiCorporationId} />
        )}
        {props.type !== "alliance" && props.esiAllianceId && (
          <AllianceAvatar size={32} id={props.esiAllianceId} />
        )}
      </div>
    </Link>
  );
}
