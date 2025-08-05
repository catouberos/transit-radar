-- +goose Up
CREATE TABLE vehicle_types (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL
);

ALTER TABLE
    vehicles
ADD
    COLUMN "type" bigint;

ALTER TABLE
    routes
ADD
    COLUMN operation_time text,
ADD
    COLUMN organization text,
ADD
    COLUMN ticketing text,
ADD
    COLUMN route_type text;

ALTER TABLE
    variants
ADD
    COLUMN description text,
ADD
    COLUMN short_name text,
ADD
    COLUMN distance real,
ADD
    COLUMN duration int,
ADD
    COLUMN start_stop_name text,
ADD
    COLUMN end_stop_name text;

-- +goose Down
DROP TABLE vehicle_types;

ALTER TABLE
    vehicles DROP COLUMN "type";

ALTER TABLE
    routes DROP COLUMN operation_time,
    DROP COLUMN organization,
    DROP COLUMN ticketing,
    DROP COLUMN route_type;

ALTER TABLE
    variants DROP COLUMN description,
    DROP COLUMN short_name,
    DROP COLUMN distance,
    DROP COLUMN duration,
    DROP COLUMN start_stop_name,
    DROP COLUMN end_stop_name;
