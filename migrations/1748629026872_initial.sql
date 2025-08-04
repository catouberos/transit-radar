-- +goose Up
CREATE TABLE vehicles (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    license_plate varchar(10) UNIQUE NOT NULL
);

CREATE TABLE routes (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number varchar(10) UNIQUE NOT NULL,
    name text NOT NULL,
    ebms_id bigint UNIQUE NULLS NOT DISTINCT,
    active boolean NOT NULL DEFAULT TRUE
);

CREATE TABLE variations (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    ebms_id bigint UNIQUE NULLS NOT DISTINCT,
    is_outbound boolean NOT NULL DEFAULT TRUE,
    route_id bigint NOT NULL REFERENCES routes(id)
);

CREATE TABLE geolocations (
    degree real NOT NULL,
    latitude real NOT NULL,
    longitude real NOT NULL,
    speed real NOT NULL,
    vehicle_id bigint NOT NULL REFERENCES vehicles(id),
    route_id bigint NOT NULL REFERENCES routes(id),
    "timestamp" timestamptz NOT NULL
) WITH (
    tsdb.hypertable,
    tsdb.partition_column = 'timestamp',
    tsdb.segmentby = 'vehicle_id',
    tsdb.orderby = 'timestamp DESC'
);

CREATE UNIQUE INDEX idx_vehicleid_timestamp ON geolocations(vehicle_id, "timestamp");

CREATE UNIQUE INDEX idx_vehicleid_routeid_timestamp ON geolocations(vehicle_id, route_id, "timestamp");

-- +goose Down
DROP TABLE geolocations;
