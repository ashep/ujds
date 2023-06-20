CREATE TABLE index
(
    id         bigserial   NOT NULL,
    name       varchar(64) NOT NULL UNIQUE,
    schema     jsonb       NOT NULL DEFAULT '{}',
    created_at timestamp   NOT NULL DEFAULT now(),
    updated_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE TABLE record_log
(
    id         bigserial   NOT NULL,
    index_id   int         NOT NULL,
    record_id  varchar(64) NOT NULL,
    data       jsonb       NOT NULL,
    created_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (index_id) REFERENCES index (id)
);

CREATE TABLE record
(
    id         varchar(64) NOT NULL,
    index_id   int         NOT NULL,
    log_id     bigint      NOT NULL,
    checksum   bytea       NOT NULL UNIQUE,
    created_at timestamp   NOT NULL DEFAULT now(),
    updated_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (id, index_id),
    FOREIGN KEY (log_id) REFERENCES record_log (id)
);

CREATE INDEX idx_record_updated_at ON record (updated_at);