import Avatar from "@/app/ui/Search/Avatar";
import { clsx } from "clsx";

export interface Props {
  esi_character_id?: number;
  esi_corporation_id: number;
  esi_alliance_id?: number;
  name: string;
  security_status?: number;
}

function resultAvatar(result: Props) {
  if (result.esi_character_id) {
    return (
      <Avatar
        alt={result.name}
        size={32}
        src={`https://images.evetech.net/characters/${result.esi_character_id}/portrait?tenant=tranquility&size=256`}
      />
    );
  }

  return null;
}

function corpAvatar(corpID: number) {
  return (
    <Avatar
      alt="corp"
      size={32}
      src={`https://images.evetech.net/corporations/${corpID}/logo?size=128`}
    />
  );
}

function allianceAvatar(allianceID: number) {
  return (
    <Avatar
      alt="alliance"
      size={32}
      src={`https://images.evetech.net/Alliance/${allianceID}_64.png`}
    />
  );
}

function secStatusColor(secStatus: number) {
  return clsx({
    "text-green-500": secStatus > 0,
    "text-red-500": secStatus < 0,
    "text-gray-500": secStatus === 0,
  });
}

export default function ResultListItem(props: Props) {
  return (
    <div className="flex justify-between">
      <div className="flex gap-1 items-start">
        {resultAvatar(props)}
        <div className="flex flex-col gap-1">
          <h2>{props.name}</h2>
          {props.security_status && (
            <h3 className={secStatusColor(props.security_status)}>
              {props.security_status}
            </h3>
          )}
        </div>
      </div>
      <div className="flex gap-1">
        {corpAvatar(props.esi_corporation_id)}
        {props.esi_alliance_id && allianceAvatar(props.esi_alliance_id)}
      </div>
    </div>
  );
}
