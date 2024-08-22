import clsx from "clsx";
import Plaque from "@/app/ui/Plaque/Plaque";
import ErrorCard from "@/app/ui/Error/Card";
import Table from "@/app/ui/KillmailParticipants/Table";
import { KillmailParticipant, ParticipantType } from "@/app/lib/definitions";
import {
  getCachedTopAttackers,
  getCachedTopVictims,
} from "@/app/actions/rivals";

function capitalizeFirstLetter(string: string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

async function getParticipants(
  id: number,
  role: "attacker" | "victim",
  type: ParticipantType
) {
  try {
    return role === "attacker"
      ? getCachedTopAttackers(id, type)
      : getCachedTopVictims(id, type);
  } catch (err) {
    throw err;
  }
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
  try {
    const data = await getParticipants(id, role, type);
    const title = capitalizeFirstLetter(role);

    let top: KillmailParticipant;
    let rest: KillmailParticipant[] = [];
    if (data.length) {
      top = data[0];
    }
    rest = data.slice(1);

    return (
      <div className="flex flex-col gap-8 justify-start items-center">
        {data?.length ? (
          <Plaque
            {...{
              title: `Top ${title}`,
              participant: top!,
              role,
            }}
          />
        ) : (
          <h2>No {`${title}s`} </h2>
        )}
        {rest.length ? (
          <>
            <div
              className={clsx(
                "divider",
                `divider-${role === "attacker" ? "warning" : "secondary"}`,
                `divider-${role === "attacker" ? "start" : "end"}`
              )}
            >
              {title}s
            </div>
            <Table participants={rest} role={role} />
          </>
        ) : null}
      </div>
    );
  } catch (err) {
    return (
      <ErrorCard
        message={`Failed to fetch ${capitalizeFirstLetter(role)}s`}
        level="warning"
      />
    );
  }
}
