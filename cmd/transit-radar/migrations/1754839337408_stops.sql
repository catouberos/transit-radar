-- +goose Up
CREATE TABLE stops (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    code varchar(10) UNIQUE NOT NULL,
    name text NOT NULL,
    type_id bigint NOT NULL,
    ebms_id bigint UNIQUE NULLS NOT DISTINCT,
    active boolean NOT NULL DEFAULT TRUE,
    latitude real NOT NULL,
    longitude real NOT NULL,
    -- address
    address_number text,
    address_street text,
    address_ward text,
    address_district text
);

CREATE INDEX idx_stop_name ON stops(name);

CREATE TABLE stop_types (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL
);

CREATE INDEX idx_stoptype_name ON stop_types(name);

CREATE TABLE variants_stops (
    variant_id bigint,
    stop_id bigint,
    order_score int,
    PRIMARY KEY (variant_id, stop_id)
);

-- +goose Down
DROP TABLE stops;

DROP TABLE stop_types;

DROP TABLE variants_stops;
