ALTER TABLE deliveries
    RENAME COLUMN type_id TO truck_id;

ALTER TABLE delivery_types
    RENAME TO trucks;
