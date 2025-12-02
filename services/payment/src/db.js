import pkg from "pg";
const { Pool } = pkg;

let pool = null;

export async function initDB() {
  if (pool) {
    return pool;
  }

  pool = new Pool({
    host: process.env.DB_HOST || "payment-service-database",
    port: process.env.DB_PORT || 5432,
    database: process.env.DB_NAME || "ridenow_payments",
    user: process.env.DB_USER || "payment_user",
    password: process.env.DB_PASSWORD || "payment_password",
  });

  // Test connection
  try {
    const client = await pool.connect();
    console.log("✅ Connecté à PostgreSQL");
    client.release();
  } catch (err) {
    console.error("❌ Erreur de connexion à PostgreSQL:", err);
    throw err;
  }

  return pool;
}

export function getDB() {
  if (!pool) {
    throw new Error("Database not initialized. Call initDB() first.");
  }
  return pool;
}
