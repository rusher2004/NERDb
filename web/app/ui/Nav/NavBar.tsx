"use client";

import { useEffect, useState } from "react";

export default function NavBar() {
  const themes = [
    { label: "Default (Photon)", value: "photon" },
    { label: "Amarr", value: "amarr" },
    { label: "Caldari", value: "caldari" },
    { label: "Gallente", value: "gallente" },
    { label: "Minmatar", value: "minmatar" },
    { label: "ORE", value: "ore" },
    { label: "SoE", value: "sisters" },
  ];

  const [chosenTheme, setChosenTheme] = useState(() => {
    const localTheme = localStorage.getItem("theme");
    return localTheme || "photon";
  });

  useEffect(() => {
    localStorage.setItem("theme", chosenTheme);
  }, [chosenTheme]);

  return (
    <nav className="navbar">
      <div className="flex-1">
        <a className="btn btn-ghost text-4xl text-primary">NERDb</a>
      </div>
      <div className="flex-none gap-2">
        <div className="form-control">
          <input
            type="text"
            placeholder="Search"
            className="input input-bordered w-24 md:w-auto"
          />
        </div>
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
            className="dropdown-content bg-base-200 rounded-box z-[1] w-52 p-2 shadow-2xl"
          >
            {themes.map((theme) => (
              <li key={theme.value}>
                <input
                  data-theme={theme.value}
                  type="radio"
                  name="theme-dropdown"
                  className="theme-controller btn btn-sm btn-block btn-ghost justify-start text-primary"
                  aria-label={theme.label}
                  value={theme.value}
                  checked={chosenTheme === theme.value}
                  onChange={() => setChosenTheme(theme.value)}
                />
              </li>
            ))}
          </ul>
        </div>
      </div>
    </nav>
  );
}
