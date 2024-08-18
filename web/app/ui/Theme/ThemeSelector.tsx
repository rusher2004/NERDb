"use client";

import Image from "next/image";
import { useEffect, useLayoutEffect, useState } from "react";

export default function ThemeSelector() {
  const themes = [
    { label: "Default (Photon)", value: "photon" },
    { label: "Amarr", value: "amarr", corpNum: 500003 },
    { label: "Caldari", value: "caldari", corpNum: 500001 },
    { label: "Gallente", value: "gallente", corpNum: 500004 },
    { label: "Minmatar", value: "minmatar", corpNum: 500002 },
    { label: "ORE", value: "ore", corpNum: 500014 },
    { label: "SoE", value: "sisters", corpNum: 500016 },
  ];
  const [chosenTheme, setChosenTheme] = useState<string | null>(null);

  useLayoutEffect(() => {
    const localTheme = localStorage.getItem("theme");

    setChosenTheme(localTheme || "photon");
  }, []);

  useEffect(() => {
    if (!chosenTheme) return;

    localStorage.setItem("theme", chosenTheme);
  }, [chosenTheme]);

  return (
    <div className="dropdown dropdown-end">
      <div tabIndex={0} role="button" className="btn m-1">
        Theme
        <svg
          width="12px"
          height="12px"
          className="inline-block h-2 w-2 fill-current opacity-60"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 2048 2048"
        >
          <path d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z"></path>
        </svg>
      </div>
      <ul
        tabIndex={0}
        className="dropdown-content bg-base-200 rounded-box z-[1] w-44 p-2 shadow-2xl"
      >
        {themes.map((theme) => (
          <li key={theme.value} className="flex">
            {theme.corpNum && (
              <Image
                alt={theme.label}
                width={32}
                height={32}
                src={`https://images.evetech.net/corporations/${theme.corpNum}/logo?size=64`}
              />
            )}
            <input
              data-theme={theme.value}
              type="radio"
              name="theme-dropdown"
              className="theme-controller btn btn-sm btn-ghost flex-grow justify-start text-primary"
              aria-label={theme.label}
              value={theme.value}
              checked={chosenTheme === theme.value}
              onChange={() => setChosenTheme(theme.value)}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}
