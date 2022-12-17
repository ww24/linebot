-- Transform geolite2.GeoLite2-City
SELECT
  country_iso_code AS country,
  city_name AS city,
  NET.IP_FROM_STRING(REGEXP_EXTRACT(network, r'(.*)/')) network,
  CAST(REGEXP_EXTRACT(network, r'/(.*)') AS INT64) mask
FROM `geolite2.GeoLite2-City-Blocks`
JOIN `geolite2.GeoLite2-City-Locations`
USING(geoname_id)
