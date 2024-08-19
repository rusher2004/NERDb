import { ParticipantType } from "@/app/lib/definitions";
import {
  getCachedTopAttackers,
  getCachedTopVictims,
} from "@/app/actions/rivals";
import KillmailParticipantCard from "@/app/ui/KillmailParticipants/KillmailParticipant";
import Plaque from "../Plaque/Plaque";

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

  return <Plaque {...{ participant: data[0], role }} />;
}
