-- +goose Up
CREATE TABLE vehicles (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    license_plate varchar(10) UNIQUE NOT NULL
);

CREATE UNIQUE INDEX idx_vehicle_licenseplate ON vehicles USING HASH (license_plate);

CREATE TABLE variations (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    ebms_id bigint UNIQUE NULLS NOT DISTINCT,
    is_outbound boolean NOT NULL DEFAULT TRUE,
    route_id bigint NOT NULL REFERENCES routes(id)
);

CREATE INDEX idx_variation_ebmsid ON variations USING HASH (ebms_id);

CREATE UNIQUE INDEX idx_variation_outbound_routeid ON variations(is_outbound, route_id);

CREATE TABLE routes (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number varchar(10) UNIQUE NOT NULL,
    name text NOT NULL,
    ebms_id bigint UNIQUE NULLS NOT DISTINCT,
    active boolean NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_route_ebmsid ON routes USING HASH (ebms_id);

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

CREATE UNIQUE INDEX idx_geolocation_vehicleid_timestamp ON geolocations(vehicle_id, "timestamp");

CREATE UNIQUE INDEX idx_geolocation_vehicleid_routeid_timestamp ON geolocations(vehicle_id, route_id, "timestamp");

-- +goose Down
DROP TABLE vehicles;

DROP TABLE variations;

DROP TABLE routes;

DROP TABLE geolocations;
