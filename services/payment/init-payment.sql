-- Database and user are created by PostgreSQL init scripts
-- This file can be used for additional setup if needed

CREATE TABLE IF NOT EXISTS payments (
  payment_id VARCHAR(255) PRIMARY KEY,
  ride_id VARCHAR(255),
  amount DECIMAL(10, 2),
  status VARCHAR(50),
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
