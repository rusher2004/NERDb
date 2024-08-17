"use server";

import postgres from "postgres";

const db = postgres(process.env.POSTGRES_URL || "", {
  // ssl: process.env.NODE_ENV === "production",
  ssl: false,
  transform: postgres.camel,
  max: 1000,
});

export default db;
