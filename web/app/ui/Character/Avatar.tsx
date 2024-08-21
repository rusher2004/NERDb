import BaseAvatar from "@/app/ui/Avatar/Avatar";

export default function Avatar({ size, id }: { size: number; id: number }) {
  return (
    <BaseAvatar
      alt="avatar"
      size={256}
      src={`https://images.evetech.net/characters/${id}/portrait?tenant=tranquility&size=256`}
    />
  );
}
