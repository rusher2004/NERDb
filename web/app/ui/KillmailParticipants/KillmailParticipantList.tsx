import { ParticipantType } from "@/app/lib/definitions";
import {
  getCachedTopAttackers,
  getCachedTopVictims,
} from "@/app/actions/rivals";
import KillmailParticipantCard from "@/app/ui/KillmailParticipants/KillmailParticipant";

export default async function KillmailParticipantList({
  id,
  type,
  side,
}: {
  id: number;
  type: ParticipantType;
  side: "attacker" | "victim";
}) {
  const data =
    side === "attacker"
      ? await getCachedTopAttackers(id, type)
      : await getCachedTopVictims(id, type);

  return data.map((participant) => (
    <KillmailParticipantCard
      key={participant.esiCharacterId}
      {...participant}
    />
  ));
}
