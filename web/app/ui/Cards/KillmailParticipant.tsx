import { getCharacter } from "@/app/actions/character";
import { KillmailParticipant } from "@/app/lib/definitions";
import CharAvatar from "@/app/ui/Character/Avatar";
import Link from "next/link";

export default async function KillmailParticipantCard(
  props: KillmailParticipant
) {
  const character = await getCharacter(props.esiCharacterId);

  return (
    <div className="stat">
      <div className="stat-figure">
        <Link href={`/character/${props.esiCharacterId}`}>
          <CharAvatar size={128} id={props.esiCharacterId} />
        </Link>
      </div>
      <div className="stat-value">{props.numberOfKills}</div>
      <div className="stat-title"># of Kills</div>
      <div className="stat-description">{character?.name}</div>
      {/* <div>{props.esiCharacterId}</div> */}
    </div>
  );
}
