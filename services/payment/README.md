# Payment Service (mock)

Endpoints:
- POST /payments/authorize  { ride_id, amount } -> 201 { payment_id, status }
- POST /payments/capture    { payment_id } -> 200 { payment_id, status }

DB: SQLite stored in /data/payments.db (container volume)

To run locally (without Docker):
1. npm install
2. NODE_ENV=development node src/app.js
