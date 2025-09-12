-- migrate:up
SELECT
  add_retention_policy ('iot_telemetry_events', INTERVAL '3 months');

SELECT
  add_retention_policy ('iot_alert_events', INTERVAL '1 year');

SELECT
  add_retention_policy ('iot_production_events', INTERVAL '1 year');

SELECT
  add_retention_policy ('iot_status_events', INTERVAL '1 year');

-- migrate:down
SELECT
  remove_retention_policy ('iot_telemetry_events');

SELECT
  remove_retention_policy ('iot_alert_events');

SELECT
  remove_retention_policy ('iot_production_events');

SELECT
  remove_retention_policy ('iot_status_events');