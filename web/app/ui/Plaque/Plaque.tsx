import clsx from "clsx";
import Link from "next/link";
import { getCharacter } from "@/app/actions/character";
import { KillmailParticipant } from "@/app/lib/definitions";

export default async function Plaque({
  participant,
  role,
}: {
  participant: KillmailParticipant;
  role: "attacker" | "victim";
}) {
  const character = await getCharacter(participant.esiCharacterId);

  return (
    <div className="flex flex-col justify-around items-center gap-1 bg-accent text-accent-content p-3 rounded-lg">
      <Link
        className="flex flex-col justify-center items-center"
        href={`/character/${participant.esiCharacterId}`}
      >
        <div className="avatar bg-base-200 px-1">
          <div className="w-32">
            <img
              src={`https://images.evetech.net/characters/${participant.esiCharacterId}/portrait?tenant=tranquility&size=128`}
              alt="avatar"
            />
          </div>
        </div>
        <div className="text-left">{character?.name}</div>
      </Link>
      <h2
        className={clsx(
          ["text-center", "text-xl", "rounded-s", "rounded-e", "w-40"],
          role === "attacker"
            ? ["bg-primary", "text-primary-content"]
            : ["bg-secondary", "text-secondary-content"]
        )}
      >
        Top {role === "attacker" ? "Attacker" : "Victim"}
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
