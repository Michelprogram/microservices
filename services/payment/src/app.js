import express from "express";
import paymentRoutes from "./routes.js";

const app = express();
app.use(express.json());

app.use("/payments", paymentRoutes);

const PORT = process.env.PORT || 8004;
app.listen(PORT, () => {
  console.log(`Payment service running on port ${PORT}`);
});
