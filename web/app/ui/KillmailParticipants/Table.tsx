import { KillmailParticipant } from "@/app/lib/definitions";
import TableRow from "@/app/ui/KillmailParticipants/TableRow";
export default function Table({
  participants,
  role,
}: {
  participants: KillmailParticipant[];
  role: "attacker" | "victim";
}) {
  return (
    <table className="table">
      <thead>
        <tr>
          <th className="text-left">Character</th>
          <th></th>
          <th className="text-right">
            # of {role === "attacker" ? "Kills" : "Deaths"}
          </th>
        </tr>
      </thead>
      <tbody>
        {participants.map((participant) => (
          <TableRow {...participant} />
        ))}
      </tbody>
    </table>
  );
}
