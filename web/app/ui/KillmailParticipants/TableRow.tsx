import Link from "next/link";
import { KillmailParticipant } from "@/app/lib/definitions";

import { getCharacter } from "@/app/actions/character";

export default async function TableRow(props: KillmailParticipant) {
  const character = await getCharacter(props.esiCharacterId);

  return (
    <tr>
      <th className="px-1">
        <Link href={`/character/${props.esiCharacterId}`}>
          <div className="avatar px-0" data-tip={character?.name}>
            <div className="w-16 rounded">
              <img
                src={`https://images.evetech.net/characters/${character?.esiCharacterId}/portrait?tenant=tranquility&size=64`}
                alt="avatar"
              />
            </div>
          </div>
        </Link>
      </th>

      <td className="text-left px-0">
        <Link href={`/character/${props.esiCharacterId}`}>
          {character?.name}
        </Link>
      </td>
      <td className="text-right px-2"> {props.numberOfKills}</td>
    </tr>
  );
}
