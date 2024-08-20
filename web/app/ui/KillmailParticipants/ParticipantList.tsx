import { KillmailParticipant, ParticipantType } from "@/app/lib/definitions";
import {
  getCachedTopAttackers,
  getCachedTopVictims,
} from "@/app/actions/rivals";
import Plaque from "@/app/ui/Plaque/Plaque";
import Table from "@/app/ui/KillmailParticipants/Table";

function capitalizeFirstLetter(string: string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

export default async function KillmailParticipantList({
  id,
  type,
  role: role,
}: {
  id: number;
  type: ParticipantType;
  role: "attacker" | "victim";
}) {
  const data =
    role === "attacker"
      ? await getCachedTopAttackers(id, type)
      : await getCachedTopVictims(id, type);

  let top: KillmailParticipant;
  if (data.length) {
    top = data[0];
  }

  const rest = data.slice(1);

  return data?.length ? (
    <div className="flex flex-col gap-8 justify-start items-center">
      <Plaque
        {...{
          title: `Top ${capitalizeFirstLetter(role)}`,
          participant: top!,
          role,
        }}
      />
      <h2 className="text-2xl">{capitalizeFirstLetter(role)}s</h2>
      <Table participants={rest} role={role} />
    </div>
  ) : (
    <h2>No {`${capitalizeFirstLetter(role)}s`} </h2>
  );
}
