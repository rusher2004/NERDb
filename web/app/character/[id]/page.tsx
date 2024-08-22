import SubjectView from "@/app/ui/Rivalries/SubjectView";

export default async function Page({ params }: { params: { id: string } }) {
  const id = parseInt(params.id);

  return <SubjectView params={{ id, type: "character" }} />;
}
