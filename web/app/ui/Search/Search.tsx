"use client";

import ResultList from "@/app/ui/Search/ResultList";

import { useEffect, useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { searchAll, SearchAllResult } from "@/app/actions/search";
import { clsx } from "clsx";

export default function Search() {
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState<SearchAllResult[]>([]);

  const debounced = useDebouncedCallback((term) => {
    setSearchTerm(term);
  }, 300);

  useEffect(() => {
    let searching = true;

    const doSearch = async (term: string) => {
      const res = await searchAll(term);

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
          placeholder="Search to find for your rivals..."
          className="grow"
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
