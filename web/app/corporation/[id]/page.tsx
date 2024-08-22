import SubjectView from "@/app/ui/Rivalries/SubjectView";

export default async function Page({ params }: { params: { id: string } }) {
  return (
    <SubjectView params={{ id: parseInt(params.id), type: "corporation" }} />
  );
}
