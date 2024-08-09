"use client";

import ResultList from "@/app/ui/Search/ResultList";

import { useEffect, useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { search, SearchResult } from "@/app/actions/search";

export default function Search() {
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);

  const debounced = useDebouncedCallback((term) => {
    setSearchTerm(term);
  }, 300);

  useEffect(() => {
    let searching = true;

    const doSearch = async (term: string) => {
      const res = await search(term);

      if (searching) {
        if (res.error) {
          console.error(res.error);
          return;
        }

        setSearchResults(res.rows || []);
      }
    };

    doSearch(searchTerm).catch(console.error);
  }, [searchTerm]);

  return (
    <div className="form-control">
      <input
        type="text"
        placeholder="Search"
        className="input input-bordered input-primary w-24 md:w-auto"
        onInput={(e) => {
          debounced((e.target as HTMLInputElement).value);
        }}
      />
      {searchResults.length ? (
        <ResultList searchResults={searchResults} />
      ) : null}
    </div>
  );
}
