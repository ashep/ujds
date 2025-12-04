ALTER table record
    ADD COLUMN touched_at TIMESTAMP;

UPDATE record
SET touched_at=r.updated_at
FROM (SELECT id, updated_at FROM record) AS r
WHERE record.id = r.id;

ALTER TABLE record
    ALTER COLUMN touched_at SET NOT NULL;

ALTER TABLE record
    ALTER COLUMN touched_at SET DEFAULT now();
