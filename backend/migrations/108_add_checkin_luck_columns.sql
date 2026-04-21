-- Add luck check-in columns to checkins table.
-- checkin_type: 'normal' (default) or 'luck'
-- bet_amount: the amount the user bet (only for luck check-in)
-- multiplier: the random multiplier applied (only for luck check-in)

ALTER TABLE checkins ADD COLUMN IF NOT EXISTS checkin_type VARCHAR(20) NOT NULL DEFAULT 'normal';
ALTER TABLE checkins ADD COLUMN IF NOT EXISTS bet_amount DECIMAL(20, 8) NOT NULL DEFAULT 0;
ALTER TABLE checkins ADD COLUMN IF NOT EXISTS multiplier DECIMAL(20, 8) NOT NULL DEFAULT 0;
