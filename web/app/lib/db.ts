import postgres from "postgres";

const sql = postgres(process.env.POSTGRES_URL || "", {
  ssl: process.env.NODE_ENV === "production" ? "prefer" : false,
  database: "nerdb",
  transform: postgres.camel,
  max: 1000,
  idle_timeout: 20,
  max_lifetime: 60 * 30,
});

export default sql;
