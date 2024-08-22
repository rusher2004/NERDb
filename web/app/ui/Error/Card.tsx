import clsx from "clsx";
import { IconAlertTriangle } from "@tabler/icons-react";

export default function Card({
  level,
  message,
}: {
  level: "info" | "warning" | "error";
  message: String;
}) {
  return (
    <div>
      <div role="alert" className={clsx("alert flex", `alert-${level}`)}>
        <div
          className="tooltip"
          data-tip="The error has been reported. Try again later."
        >
          <IconAlertTriangle />
        </div>
        <span>{message}</span>
      </div>
    </div>
  );
}
