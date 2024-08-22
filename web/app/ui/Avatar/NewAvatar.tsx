import {
  Alliance,
  Character,
  Corporation,
  Participant,
  ParticipantType,
} from "@/app/lib/definitions";

export default function Avatar({
  subject,
  type,
}: {
  subject: Participant;
  type: ParticipantType;
}) {
  const srcURL = (subject: Participant) => {
    switch (type) {
      case "character":
        const char = subject as Character;
        return `https://images.evetech.net/characters/${char.esiCharacterId}/portrait?size=256`;
      case "corporation":
        const corp = subject as Corporation;
        return `https://images.evetech.net/corporations/${corp.esiCorporationId}/logo?size=256`;
      case "alliance":
        const alliance = subject as Alliance;
        return `https://images.evetech.net/Alliance/${alliance.esiAllianceId}_64.png`;
      default:
        return "";
    }
  };

  return (
    <div className="avatar">
      <div className="max-w-48 max-h-48">
        <img src={srcURL(subject)} alt={subject.name} />
      </div>
    </div>
  );
}
