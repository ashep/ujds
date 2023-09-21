CREATE TABLE index
(
    id         BIGSERIAL   NOT NULL,
    parent_id  BIGINT      NULL,
    name       VARCHAR(64) NOT NULL UNIQUE,
    schema     JSONB       NOT NULL DEFAULT '{}',
    created_at TIMESTAMP   NOT NULL DEFAULT now(),
    updated_at TIMESTAMP   NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES index (id)
);

CREATE TABLE record_log
(
    id         BIGSERIAL   NOT NULL,
    index_id   INT         NOT NULL,
    record_id  VARCHAR(64) NOT NULL,
    data       JSONB       NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (index_id) REFERENCES index (id)
);

CREATE TABLE record
(
    id         VARCHAR(64) NOT NULL,
    index_id   INT         NOT NULL,
    log_id     BIGINT      NOT NULL,
    checksum   BYTEA       NOT NULL UNIQUE,
    created_at TIMESTAMP   NOT NULL DEFAULT now(),
    updated_at TIMESTAMP   NOT NULL DEFAULT now(),

    PRIMARY KEY (id, index_id),
    FOREIGN KEY (log_id) REFERENCES record_log (id)
);

CREATE INDEX idx_record_updated_at ON record (updated_at);