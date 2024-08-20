"use client";

import Avatar from "@/app/ui/Corporation/Avatar";
import { useEffect, useLayoutEffect, useState } from "react";

interface Theme {
  label: string;
  value: string;
  corpNum: number;
}

export default function ThemeSelector() {
  const themes: Theme[] = [
    { label: "Photon", value: "photon", corpNum: 109299958 },
    { label: "Amarr", value: "amarr", corpNum: 500003 },
    { label: "Caldari", value: "caldari", corpNum: 500001 },
    { label: "Gallente", value: "gallente", corpNum: 500004 },
    { label: "Minmatar", value: "minmatar", corpNum: 500002 },
    { label: "ORE", value: "ore", corpNum: 500014 },
    { label: "SoE", value: "sisters", corpNum: 500016 },
  ];
  const [chosenTheme, setChosenTheme] = useState<Theme | null>(null);

  useLayoutEffect(() => {
    const localTheme = localStorage.getItem("theme");

    const photon = themes[0];
    const found = themes.find((theme) => theme.value === localTheme);

    setChosenTheme(found || photon);
  }, []);

  useEffect(() => {
    if (!chosenTheme) return;

    localStorage.setItem("theme", chosenTheme.value);
  }, [chosenTheme]);

  function handleThemeChange(event: React.ChangeEvent<HTMLInputElement>) {
    const theme = themes.find((theme) => theme.value === event.target.value);

    if (theme) setChosenTheme(theme);
  }

  return (
    <div className="dropdown dropdown-end">
      <div tabIndex={0} role="button" className="btn m-1">
        {chosenTheme && <Avatar size={32} id={chosenTheme.corpNum} />}
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
          <li key={theme.value} className="flex items-center">
            {theme.corpNum && <Avatar size={32} id={theme.corpNum} />}
            <input
              data-theme={theme.value}
              type="radio"
              name="theme-dropdown"
              className="theme-controller btn btn-sm btn-ghost flex-grow justify-start text-primary text-sm"
              aria-label={theme.label}
              value={theme.value}
              checked={chosenTheme?.value === theme.value}
              onChange={handleThemeChange}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}
