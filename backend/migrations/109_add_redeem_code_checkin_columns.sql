-- Add multiplier and bet_amount columns to redeem_codes table for checkin luck records.

ALTER TABLE redeem_codes ADD COLUMN IF NOT EXISTS multiplier DECIMAL(20, 8) NOT NULL DEFAULT 0;
ALTER TABLE redeem_codes ADD COLUMN IF NOT EXISTS bet_amount DECIMAL(20, 8) NOT NULL DEFAULT 0;
