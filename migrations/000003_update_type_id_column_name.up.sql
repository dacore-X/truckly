ALTER TABLE deliveries
    RENAME COLUMN truck_id TO type_id;

ALTER TABLE trucks
    RENAME TO delivery_types;