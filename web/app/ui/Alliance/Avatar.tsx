import BaseAvatar from "@/app/ui/Avatar/Avatar";

export default function Avatar({ size, id }: { size: number; id: number }) {
  return (
    <BaseAvatar
      alt="avatar"
      size={128}
      src={`https://images.evetech.net/Alliance/${id}_64.png`}
    />
  );
}
