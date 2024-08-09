import ResultListItem from "@/app/ui/Search/ResultListItem";
import { SearchResult } from "@/app/actions/search";

export interface Props {
  searchResults: SearchResult[];
}

export default function ResultList({ searchResults }: Props) {
  return (
    <ul className="menu flex-nowrap p-4 shadow-lg bg-base-100 rounded-box h-96 overflow-y-auto">
      {searchResults.map((result) => (
        <li key={result.esi_character_id}>
          <ResultListItem {...result} />
        </li>
      ))}
    </ul>
  );
}
