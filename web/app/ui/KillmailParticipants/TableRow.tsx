import { KillmailParticipant } from "@/app/lib/definitions";

import { getCharacter } from "@/app/actions/character";

export default async function TableRow(props: KillmailParticipant) {
  const character = await getCharacter(props.esiCharacterId);

  return (
    <tr>
      <th>
        <div className="avatar px-1" data-tip={character?.name}>
          <div className="w-16 rounded">
            <img
              src={`https://images.evetech.net/characters/${character?.esiCharacterId}/portrait?tenant=tranquility&size=64`}
              alt="avatar"
            />
          </div>
        </div>
      </th>
      <td className="text-left">{character?.name}</td>
      <td className="text-right"> {props.numberOfKills}</td>
    </tr>
  );
}
