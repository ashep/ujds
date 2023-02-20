CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE item
(
    id      UUID    NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    version INTEGER NOT NULL,
    type    VARCHAR NOT NULL
);

CREATE TABLE item_version
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    item_id UUID      NOT NULL,
    data    JSONB     NOT NULL DEFAULT '{}',
    time    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX item_version_idx1 ON item_version (id, item_id);

CREATE TABLE item_log
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    item_id UUID      NOT NULL,
    version BIGSERIAL NOT NULL,
    message VARCHAR   NOT NULL,
    time    TIMESTAMP NOT NULL DEFAULT NOW()
);

