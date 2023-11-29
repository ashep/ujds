ALTER TABLE record
    ADD COLUMN data JSONB;

UPDATE record
SET data=log.data
FROM (SELECT index_id, record_id, data
      FROM record_log
      WHERE id IN (SELECT max(id) FROM record_log GROUP BY record_id)) log
WHERE record.index_id = log.index_id
  AND record.id = log.record_id;

ALTER TABLE record
    ALTER COLUMN data SET NOT NULL;
