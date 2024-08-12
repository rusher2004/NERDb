"use client";

import ResultList from "@/app/ui/Search/ResultList";

import { useEffect, useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { search, SearchResult } from "@/app/actions/search";
import { clsx } from "clsx";

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
    <div
      tabIndex={0}
      // className="form-control"
      className={clsx({
        "form-control": true,
        collapse: true,
        "collapse-open": searchResults.length,
        "collapse-close": !searchResults.length,
        "rounded-none": true,
      })}
    >
      <label className="input input-bordered input-primary collapse-title flex items-center gap-2">
        <input
          type="text"
          placeholder="Search your character to find your rival"
          className="grow"
          // className="input input-bordered input-primary w-24 md:w-auto"
          onInput={(e) => {
            debounced((e.target as HTMLInputElement).value);
          }}
        />
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 16 16"
          fill="currentColor"
          className="h-4 w-4 opacity-70"
        >
          <path
            fillRule="evenodd"
            d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z"
            clipRule="evenodd"
          />
        </svg>
      </label>
      <ResultList searchResults={searchResults} />
    </div>
  );
}
