import { getCharacter } from "@/app/actions/character";
import { KillmailParticipant } from "@/app/lib/definitions";
import CharAvatar from "@/app/ui/Character/Avatar";

export default async function KillmailParticipantCard(
  props: KillmailParticipant
) {
  const character = await getCharacter(props.esiCharacterId);

  return (
    <div className="stat">
      <div className="stat-figure">
        <CharAvatar size={64} id={props.esiCharacterId} />
      </div>
      <div className="stat-value">{props.numberOfKills}</div>
      <div className="stat-title"># of Kills</div>
      <div className="stat-description">{character?.name}</div>
      {/* <div>{props.esiCharacterId}</div> */}
    </div>
  );
}
