import BaseAvatar from "@/app/ui/Avatar/Avatar";

export default function Avatar({ size, id }: { size: number; id: number }) {
  return (
    <div className="avatar">
      <div className="max-w-16 max-h-16">
        <img
          src={`https://images.evetech.net/corporations/${id}/logo?size=256`}
          alt="corporation logo"
        />
      </div>
    </div>
    // <BaseAvatar
    //   alt="avatar"
    //   size={128}
    //   src={`https://images.evetech.net/corporations/${id}/logo?size=128`}
    // />
  );
}
