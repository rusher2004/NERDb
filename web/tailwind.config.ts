import type { Config } from "tailwindcss";
import daisyui from "daisyui";

const eveThemeDefaults = {
  "base-100": "#111111",
  neutral: "#C8C8C1",
};

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
    },
  },
  daisyui: {
    themes: [
      {
        photon: {
          ...eveThemeDefaults,
          primary: "#58A7BF",
          secondary: "#373837",
          accent: "#58A7BF",
          "base-200": "#060A0C",
          warning: "#F39058",
        },
        amarr: {
          ...eveThemeDefaults,
          primary: "#A38264",
          secondary: "#50412B",
          accent: "FFEE93",
          "base-200": "#07030D",
          warning: "#F39058",
        },
        caldari: {
          ...eveThemeDefaults,
          primary: "#9AD2E3",
          secondary: "#282827",
          accent: "#48C7BF",
          "base-200": "#030E0E",
          warning: "#DC8C2E",
        },
        gallente: {
          ...eveThemeDefaults,
          primary: "#58BF9A",
          secondary: "#344238",
          accent: "#6DB09E",
          "base-200": "#0A090F",
          warning: "#F39058",
        },
        minmatar: {
          ...eveThemeDefaults,
          primary: "#D05C3B",
          secondary: "#362A25",
          accent: "#9D452D",
          "base-200": "#030909",
          warning: "#F39058",
        },
        ore: {
          ...eveThemeDefaults,
          primary: "#DDB825",
          secondary: "#394961",
          accent: "#55999C",
          "base-200": "#030805",
          warning: "#DE8B78",
        },
        sisters: {
          ...eveThemeDefaults,
          primary: "#E55252",
          secondary: "#373837",
          accent: "#A1DDE0",
          "base-200": "#070D11",
          warning: "#C0B337",
        },
      },
    ],
  },
  plugins: [daisyui, require("@tailwindcss/typography")],
};
export default config;

// amarr #A38264,#FFEE93,#07030D,#F39058
// caldari #9AD2E3,#48C7BF,#030E0E,#DC8C2E
// gallente #58BF9A,#6DB09E,#0A090F,#F39058
// minmatar #D05C3B,#9D452D,#030909,#F39058
// ore #DDB825,#55999C,#030805,#DE8B78
// photon #58A7BF,#58A7BF,#060A0C,#F39058
// sisters #E55252,#A1DDE0,#070D11,#C0B337
