-- with_geolocation function
-- CREATE OR REPLACE TABLE FUNCTION ${dataset}.with_geolocation(since TIMESTAMP, until TIMESTAMP) AS
WITH
  access_logs AS (SELECT *
    FROM `${project}.${dataset}.access_log`
    WHERE `timestamp` BETWEEN since AND until),
  geolocations AS (SELECT *
    FROM `${project}.geolite2.GeoLite2_City_*`
    WHERE _TABLE_SUFFIX = FORMAT_DATE('%Y%m%d', DATE(since)))
SELECT * FROM access_logs
LEFT JOIN (
  WITH ips AS (SELECT DISTINCT ip FROM access_logs)
  -- IPv4 address => country, city
  SELECT ip, country, city FROM (
    SELECT NET.IP_TRUNC(NET.SAFE_IP_FROM_STRING(ip), mask) network, *
    FROM ips, UNNEST(GENERATE_ARRAY(8,32)) mask
    WHERE ip LIKE '%.%'
  )
  JOIN geolocations USING (network, mask)
  UNION ALL
  -- IPv6 address => country, city
  SELECT ip, country, city FROM (
    SELECT NET.IP_TRUNC(NET.SAFE_IP_FROM_STRING(ip), mask) network, *
    FROM ips, UNNEST(GENERATE_ARRAY(19,64)) mask
    WHERE ip LIKE '%:%'
  )
  JOIN geolocations USING (network, mask)
) USING (ip)
