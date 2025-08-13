-- migrate:up
CREATE TABLE
    IF NOT EXISTS iot_alert_events (
        time TIMESTAMPTZ NOT NULL,
        device_id VARCHAR(50) NOT NULL,
        alert_type VARCHAR(50) NOT NULL,
        severity VARCHAR(20) NOT NULL,
        message TEXT NOT NULL,
        current_value NUMERIC(6, 2)
    )
WITH
    (
        tsdb.hypertable,
        tsdb.partition_column = 'time',
        tsdb.segmentby = 'device_id',
        tsdb.orderby = 'time DESC'
    );

-- migrate:down
DROP TABLE IF EXISTS iot_alert_events CASCADE;