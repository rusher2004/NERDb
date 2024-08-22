import clsx from "clsx";
import Link from "next/link";
import { getCharacter } from "@/app/actions/character";
import { KillmailParticipant } from "@/app/lib/definitions";

export default async function Plaque({
  title,
  participant,
  role,
}: {
  title: string;
  participant: KillmailParticipant;
  role: "attacker" | "victim";
}) {
  const character = await getCharacter(participant.esiCharacterId);

  return (
    <div className="flex flex-col justify-around items-center gap-1 bg-accent text-accent-content p-3 rounded-lg w-48">
      <Link
        className="flex flex-col justify-center items-center w-full"
        href={`/character/${participant.esiCharacterId}`}
      >
        <div className="avatar px-1" data-tip={character?.name}>
          <div className="w-32 rounded">
            <img
              src={`https://images.evetech.net/characters/${participant.esiCharacterId}/portrait?tenant=tranquility&size=128`}
              alt="avatar"
            />
          </div>
        </div>
        <div
          className="text-center truncate w-full tooltip tooltip-open"
          data-tip={character?.name}
        >
          {character?.name}
        </div>
      </Link>
      <h2
        className={clsx(
          "text-center text-xl rounded-s rounded-e w-40",
          role === "attacker"
            ? ["bg-warning text-warning-content"]
            : ["bg-secondary text-secondary-content"]
        )}
      >
        {title}
      </h2>
      <div className="join w-40">
        <div className="join-item flex justify-between items-center w-full">
          <div># of {role === "attacker" ? "kills" : "deaths"}:</div>
          <div className="">{participant.numberOfKills}</div>
        </div>
      </div>
    </div>
  );
}
