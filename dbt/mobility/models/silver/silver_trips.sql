-- Reads from Delta Bronze (delta catalog) and selects useful fields
select
  trip_id, rider_id, driver_id, market, status, ts, duration_sec, distance_km, fare, event_date
from delta.default."mobility-delta/bronze/trips"
where event_date >= current_date - interval '7' day
