CREATE TABLE schema
(
    id         int         NOT NULL UNIQUE,
    name       varchar(64) NOT NULL UNIQUE,
    data       jsonb       NOT NULL DEFAULT '{}',
    created_at timestamp   NOT NULL DEFAULT now(),
    updated_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE SEQUENCE schema_id OWNED BY schema.id;

CREATE TABLE record_log
(
    id         bigint      NOT NULL UNIQUE,
    schema_id  int         NOT NULL,
    record_id  varchar(64) NOT NULL,
    checksum   char(64)    NOT NULL UNIQUE,
    data       jsonb       NOT NULL DEFAULT '{}',
    created_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (schema_id) REFERENCES schema (id)
);

CREATE SEQUENCE record_log_id OWNED BY record_log.id;

CREATE TABLE record
(
    id         varchar(64) NOT NULL,
    schema_id  int         NOT NULL,
    version    bigint      NOT NULL,
    created_at timestamp   NOT NULL DEFAULT now(),
    updated_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (schema_id) REFERENCES schema (id),
    FOREIGN KEY (version) REFERENCES record_log (id)
);

