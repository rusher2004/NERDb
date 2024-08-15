import Avatar from "@/app/ui/Avatar/Avatar";
import AllianceAvatar from "@/app/ui/Alliance/Avatar";
import CharAvatar from "@/app/ui/Character/Avatar";
import CorpAvatar from "@/app/ui/Corporation/Avatar";
import { clsx } from "clsx";

export interface Props {
  id: number;
  name: string;
  esi_corporation_id?: number;
  esi_alliance_id?: number;
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
    <div className="flex justify-between">
      <div className="flex gap-1 items-start">
        {resultAvatar(props)}
        <div className="flex flex-col gap-1">
          <h2>{props.name}</h2>
        </div>
      </div>
      <div className="flex gap-1">
        {props.esi_corporation_id && props.type === "character" && (
          <CorpAvatar size={32} id={props.esi_corporation_id} />
        )}
        {props.type !== "alliance" && props.esi_alliance_id && (
          <AllianceAvatar size={32} id={props.esi_alliance_id} />
        )}
      </div>
    </div>
  );
}
