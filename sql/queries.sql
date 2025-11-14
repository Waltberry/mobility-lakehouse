-- Show catalogs
SHOW CATALOGS;

-- Delta bronze table listing (Trino sees Delta paths as tables by path expression)
SELECT * FROM delta.default."mobility-delta/bronze/trips" LIMIT 5;

-- Partition pruning example
SELECT market, count(*) 
FROM delta.default."mobility-delta/bronze/trips"
WHERE event_date >= date '2025-01-01' AND market='CA-AB'
GROUP BY market;

-- Gold KPI query
SELECT * FROM gold.gold_market_hourly_kpis ORDER BY hour DESC LIMIT 20;
