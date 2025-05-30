-- +goose Up
CREATE TABLE geolocations (
    degree float4,
    latitude float4,
    longitude float4,
    speed float4,
    vehicle_id int,
    route_id int,
    "timestamp" timestamptz
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
