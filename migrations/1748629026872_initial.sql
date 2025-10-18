-- +goose Up
CREATE EXTENSION postgis;

CREATE TABLE vehicles (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    license_plate varchar(25) UNIQUE NOT NULL,
    "type" bigint
);

CREATE UNIQUE INDEX idx_vehicle_licenseplate ON vehicles(license_plate);

CREATE TABLE vehicle_types (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL
);

---- ROUTES ----
-- see `route_type` of https://gtfs.org/documentation/schedule/reference/#routestxt
CREATE TYPE route_type AS ENUM (
    'tram',
    'metro',
    'rail',
    'bus',
    'ferry',
    'cable_tram',
    'cable_car',
    'funicular',
    'trolleybus',
    'monorail'
);

CREATE TABLE routes (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number varchar(10) UNIQUE NOT NULL,
    name text NOT NULL,
    short_name text,
    description text,
    "type" route_type NOT NULL,
    color VARCHAR(7),
    agency_id bigint NOT NULL REFERENCES agencies,
    active boolean NOT NULL DEFAULT TRUE,
    attributes jsonb
);

ALTER TABLE
    routes
ADD
    CONSTRAINT color_hex_constraint CHECK (
        color IS NULL
        OR color ~* '^#[a-f0-9]{6}$'
    );

CREATE INDEX idx_route_attributes ON routes USING GIN (attributes);

---- VARIANTS ----
CREATE TABLE variants (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    route_id bigint NOT NULL,
    name text NOT NULL,
    short_name text,
    description text NOT NULL,
    distance real NOT NULL,
    direction int NOT NULL DEFAULT 0,
    duration int,
    attributes jsonb
);

CREATE INDEX idx_variant_attributes ON variants USING GIN (attributes);

CREATE UNIQUE INDEX idx_variant_direction_routeid ON variants(direction, route_id);

---- GEOLOCATIONS ----
CREATE TABLE geolocations (
    degree real NOT NULL,
    location GEOGRAPHY(POINT),
    speed real NOT NULL,
    vehicle_id BIGINT NOT NULL,
    variant_id BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
) WITH (
    tsdb.hypertable,
    tsdb.partition_column = 'timestamp',
    tsdb.segmentby = 'vehicle_id',
    tsdb.orderby = 'timestamp DESC'
);

CREATE UNIQUE INDEX idx_geolocation_vehicleid_timestamp ON geolocations(vehicle_id, "timestamp");

CREATE UNIQUE INDEX idx_geolocation_vehicleid_variantid_timestamp ON geolocations(vehicle_id, variant_id, "timestamp");

CREATE TABLE agencies (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL
);

---- STOPS ----
-- check `location_type` of https://gtfs.org/documentation/schedule/reference/#stopstxt
CREATE TYPE stop_type AS ENUM (
    'stop_platform',
    'station',
    'entrance_exit',
    'generic_node',
    'boarding_area'
);

CREATE TABLE stops (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    parent_id bigint REFERENCES stops,
    code text NOT NULL,
    name text NOT NULL,
    "type" stop_type NOT NULL,
    active boolean NOT NULL DEFAULT TRUE,
    location GEOGRAPHY(POINT),
    attributes jsonb
);

CREATE INDEX idx_stop_name ON stops(name);

CREATE INDEX idx_stop_attributes ON stops USING GIN (attributes);

---- VARIANTS STOPS ----
CREATE TABLE variants_stops (
    variant_id bigint,
    stop_id bigint,
    order_score int,
    PRIMARY KEY (variant_id, stop_id, order_score)
);

-- +goose Down
DROP TABLE stops;

DROP TABLE stop_types;

DROP TABLE variants_stops;

DROP TABLE geolocations;

DROP TABLE agencies;

DROP TABLE variants;

DROP TABLE routes;

DROP TABLE vehicle_types;

DROP TABLE vehicles;
