-- Snapshot geolite2.GeoLite2_City_YYYYMMDD
DECLARE snapshot_name STRING;
DECLARE expiration TIMESTAMP;
DECLARE query STRING;

SET expiration = TIMESTAMP_ADD(@run_time, INTERVAL 365 DAY);
SET snapshot_name = CONCAT(
  "`geolite2.GeoLite2_City_",
  FORMAT_DATE('%Y%m%d', EXTRACT(DATE FROM @run_time AT TIME ZONE "Asia/Tokyo")),
  "`"
);

SET query = CONCAT(
  "CREATE SNAPSHOT TABLE IF NOT EXISTS ",
  snapshot_name,
  " CLONE `geolite2.GeoLite2-City` OPTIONS(expiration_timestamp = TIMESTAMP '",
  expiration,
  "');"
);

EXECUTE IMMEDIATE query;
