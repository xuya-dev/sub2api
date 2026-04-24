-- Recalculate check-in streaks from the actual check-in dates.
-- This repairs legacy rows whose stored streak_days drifted away from the
-- real consecutive-day sequence.

WITH ordered AS (
    SELECT
        id,
        user_id,
        checkin_date,
        CASE
            WHEN LAG(checkin_date) OVER (PARTITION BY user_id ORDER BY checkin_date, id)
                = checkin_date - INTERVAL '1 day' THEN 0
            ELSE 1
        END AS starts_new_streak
    FROM checkins
),
grouped AS (
    SELECT
        id,
        user_id,
        checkin_date,
        SUM(starts_new_streak) OVER (
            PARTITION BY user_id
            ORDER BY checkin_date, id
            ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
        ) AS streak_group
    FROM ordered
),
recalculated AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, streak_group
            ORDER BY checkin_date, id
        ) AS streak_days
    FROM grouped
)
UPDATE checkins AS c
SET streak_days = r.streak_days
FROM recalculated AS r
WHERE c.id = r.id
  AND c.streak_days IS DISTINCT FROM r.streak_days;
