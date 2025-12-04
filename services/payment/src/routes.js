import express from "express";
import { getDB } from "./db.js";
import { v4 as uuidv4 } from "uuid";

const router = express.Router();

router.post("/authorize", async (req, res) => {
  try {
    const { ride_id, amount } = req.body;

    console.log("[PAYMENT] Authorizing payment for ride:", ride_id, "amount:", amount);

    if (!ride_id || !amount || Number(amount) <= 0) {
      return res.status(400).json({ error: "Invalid ride_id or amount" });
    }

    const payment_id = "P-" + uuidv4();
    const db = getDB();

    await db.query(
      "INSERT INTO payments (payment_id, ride_id, amount, status) VALUES ($1, $2, $3, $4)",
      [payment_id, ride_id, amount, "AUTHORIZED"]
    );

    console.log(
      `[PAYMENT] Authorized payment ${payment_id} for ride ${ride_id} amount ${amount}`
    );

    return res.status(201).json({ payment_id, status: "AUTHORIZED" });
  } catch (err) {
    console.error("[PAYMENT][ERROR] authorize:", err);
    return res.status(500).json({ error: "Internal server error" });
  }
});

router.post("/capture", async (req, res) => {
  try {
    const { payment_id } = req.body;
    if (!payment_id) {
      return res.status(400).json({ error: "payment_id required" });
    }

    const db = getDB();
    const result = await db.query(
      "SELECT * FROM payments WHERE payment_id = $1",
      [payment_id]
    );

    const payment = result.rows[0];

    if (!payment) {
      return res.status(404).json({ error: "Payment not found" });
    }

    if (payment.status === "CAPTURED") {
      return res
        .status(409)
        .json({ error: "Payment already captured", payment_id });
    }

    await db.query("UPDATE payments SET status = $1 WHERE payment_id = $2", [
      "CAPTURED",
      payment_id,
    ]);

    console.log(`[PAYMENT] Captured payment ${payment_id}`);

    return res.json({ payment_id, status: "CAPTURED" });
  } catch (err) {
    console.error("[PAYMENT][ERROR] capture:", err);
    return res.status(500).json({ error: "Internal server error" });
  }
});

export default router;
