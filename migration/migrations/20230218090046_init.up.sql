CREATE TABLE schema
(
    id   int         NOT NULL,
    name varchar(64) NOT NULL UNIQUE,
    data jsonb       NOT NULL DEFAULT '{}',

    PRIMARY KEY (name)
);

CREATE SEQUENCE schema_id OWNED BY record.id;

CREATE TABLE record_log
(
    id         varchar(64) NOT NULL,
    schema_id  int         NOT NULL,
    version    bigint      NOT NULL UNIQUE,
    checksum   char(64)    NOT NULL UNIQUE,
    data       jsonb       NOT NULL DEFAULT '{}',
    created_at timestamp   NOT NULL DEFAULT now(),

    PRIMARY KEY (version),
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
    FOREIGN KEY (version) REFERENCES record_log (version)
);

