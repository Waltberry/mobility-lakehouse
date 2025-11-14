with base as (
  select date_trunc('hour', ts) as hour, market, 
         count(*) as trips,
         avg(duration_sec) as avg_duration_sec,
         avg(distance_km) as avg_distance_km,
         avg(fare) as avg_fare
  from {{ ref('silver_trips') }}
  group by 1,2
)
select * from base
