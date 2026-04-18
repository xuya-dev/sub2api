-- Final cleanup: convert all registration bonus records from invitation to registration type
-- IDs 12,13,14,15 were created by backfill with random hex codes
-- ID 29 was created by code (before registration type fix was deployed)

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '10min';

UPDATE redeem_codes
SET type = 'registration'
WHERE type = 'invitation'
  AND value = 5
  AND status = 'used'
  AND used_by IN (2, 3, 5, 6, 9);
