import sqlite3 from "sqlite3";
import { open } from "sqlite";

export async function initDB(dbPath = "./payments.db") {
  const db = await open({
    filename: dbPath,
    driver: sqlite3.Database
  });

  await db.exec(`
    CREATE TABLE IF NOT EXISTS payments (
      payment_id TEXT PRIMARY KEY,
      ride_id TEXT,
      amount REAL,
      status TEXT,
      timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );
  `);

  return db;
}
