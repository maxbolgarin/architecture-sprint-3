CREATE TABLE IF NOT EXISTS heating_systems (
    id BIGSERIAL PRIMARY KEY,
    is_on BOOLEAN NOT NULL,
    target_temperature DOUBLE PRECISION NOT NULL,
    current_temperature DOUBLE PRECISION NOT NULL
);

CREATE TABLE IF NOT EXISTS temperature_sensors (
    id BIGSERIAL PRIMARY KEY,
    current_temperature DOUBLE PRECISION NOT NULL,
    last_updated TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS commands (
    command_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    command_type_id TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    send_time TIMESTAMPTZ,
    command_type TEXT NOT NULL,
    code TEXT NOT NULL,
    PRIMARY KEY (command_id)
);

CREATE TABLE IF NOT EXISTS telemetry (
    device_id TEXT NOT NULL,
    metric_id TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    value INT NOT NULL,
    PRIMARY KEY (device_id, metric_id, create_time)
);
