-- +goose Up
CREATE TABLE vehicles (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number_plate varchar(10) NOT NULL UNIQUE
);

CREATE TABLE routes (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number int NOT NULL UNIQUE,
    name text NOT NULL
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
