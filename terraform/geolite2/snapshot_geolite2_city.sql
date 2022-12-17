-- Snapshot geolite2.GeoLite2-City
DECLARE snapshot_name STRING;
DECLARE expiration TIMESTAMP;
DECLARE query STRING;

SET expiration = TIMESTAMP_ADD(@run_time, INTERVAL 365 DAY);
SET snapshot_name = CONCAT("`geolite2.GeoLite2_City_", FORMAT_DATETIME('%Y%m%d', @run_date), "`");

SET query = CONCAT(
  "CREATE SNAPSHOT TABLE IF NOT EXISTS ",
  snapshot_name,
  " CLONE `geolite2.GeoLite2_City` OPTIONS(expiration_timestamp = TIMESTAMP '",
  expiration,
  "');"
);

EXECUTE IMMEDIATE query;
