"use server";

import postgres from "postgres";

const sql = postgres(process.env.POSTGRES_URL || "", {
  ssl: process.env.NODE_ENV === "production" ? "prefer" : false,
  transform: postgres.camel,
  max: 1000,
});

export default sql;
