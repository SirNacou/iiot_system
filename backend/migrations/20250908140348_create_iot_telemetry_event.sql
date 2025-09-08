-- migrate:up
CREATE TABLE
    IF NOT EXISTS iot_telemetry_events (
        time TIMESTAMPTZ NOT NULL,
        device_id VARCHAR(50) NOT NULL,
        temperature_celcius NUMERIC(5,2) NOT NULL,
        humidity_percent NUMERIC(5,2) NOT NULL,
        vibration_hz NUMERIC(5,2) NOT NULL,
        motor_rpm INTEGER NOT NULL,
        current_amps NUMERIC(5,2) NOT NULL,
        machine_status VARCHAR(20) NOT NULL,
        error_code VARCHAR(50)
    )
WITH
    (
        tsdb.hypertable,
        tsdb.partition_column = 'time',
        tsdb.segmentby = 'device_id',
        tsdb.orderby = 'time DESC'
    );

-- migrate:down
DROP TABLE IF EXISTS iot_telemetry_events CASCADE;