-- migrate:up
CREATE TABLE
    IF NOT EXISTS iot_status_events (
        time TIMESTAMPTZ NOT NULL,
        device_id VARCHAR(50) NOT NULL,
        old_status VARCHAR(20) NOT NULL,
        new_status VARCHAR(20) NOT NULL,
        reason VARCHAR(50) NOT NULL
    )
WITH
    (
        tsdb.hypertable,
        tsdb.partition_column = 'time',
        tsdb.segmentby = 'device_id',
        tsdb.orderby = 'time DESC'
    );

-- migrate:down
DROP TABLE IF EXISTS iot_status_events CASCADE;