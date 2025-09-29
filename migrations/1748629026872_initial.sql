-- +goose Up
CREATE TABLE vehicles (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    license_plate varchar(25) UNIQUE NOT NULL
);

CREATE UNIQUE INDEX idx_vehicle_licenseplate ON vehicles(license_plate);

CREATE TABLE routes (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number varchar(10) UNIQUE NOT NULL,
    name text NOT NULL,
    ebms_id bigint UNIQUE NULLS NOT DISTINCT,
    active boolean NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_route_ebmsid ON routes(ebms_id);

CREATE TABLE variants (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    ebms_id bigint,
    is_outbound boolean NOT NULL DEFAULT TRUE,
    route_id bigint NOT NULL
);

CREATE INDEX idx_variant_ebmsid ON variants(ebms_id);

CREATE UNIQUE INDEX idx_variant_outbound_routeid ON variants(is_outbound, route_id);

CREATE TABLE geolocations (
    degree real NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    speed real NOT NULL,
    vehicle_id bigint NOT NULL,
    variant_id bigint NOT NULL,
    "timestamp" timestamptz NOT NULL
) WITH (
    tsdb.hypertable,
    tsdb.partition_column = 'timestamp',
    tsdb.segmentby = 'vehicle_id',
    tsdb.orderby = 'timestamp DESC'
);

CREATE UNIQUE INDEX idx_geolocation_vehicleid_timestamp ON geolocations(vehicle_id, "timestamp");

CREATE UNIQUE INDEX idx_geolocation_vehicleid_variantid_timestamp ON geolocations(vehicle_id, variant_id, "timestamp");

-- +goose Down
DROP TABLE geolocations;

DROP TABLE variants;

DROP TABLE routes;

DROP TABLE vehicles;
