-- migrate:up
CREATE TABLE
    IF NOT EXISTS iot_production_events (
        time TIMESTAMPTZ NOT NULL,
        device_id VARCHAR(50) NOT NULL,
        production_type VARCHAR(50) NOT NULL,
        product_sku VARCHAR(50) NOT NULL,
        unit_count INTEGER NOT NULL,
        batch_id VARCHAR(50) NOT NULL,
        quality_status VARCHAR(50) NOT NULL
    )
WITH
    (
        tsdb.hypertable,
        tsdb.partition_column = 'time',
        tsdb.segmentby = 'device_id',
        tsdb.orderby = 'time DESC'
    );

-- migrate:down
DROP TABLE IF EXISTS iot_production_events CASCADE;