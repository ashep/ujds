CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE auth_entity
(
    id          uuid  NOT NULL DEFAULT uuid_generate_v4(),
    secret      bytea NOT NULL,
    permissions jsonb NOT NULL DEFAULT '{}',
    note        varchar,

    PRIMARY KEY (id)
);

CREATE TABLE auth_token
(
    id             SERIAL             NOT NULL,
    auth_entity_id uuid               NOT NULL,
    token          varchar(64) UNIQUE NOT NULL DEFAULT md5(random()::text) || md5(random()::text),
    created_at     timestamp          NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    FOREIGN KEY (auth_entity_id) REFERENCES auth_entity (id)
);

CREATE TABLE item_type
(
    id     serial  NOT NULL,
    name   varchar NOT NULL,
    schema jsonb   NOT NULL DEFAULT '{}',

    PRIMARY KEY (id)
);

CREATE TABLE item
(
    id             uuid      NOT NULL DEFAULT uuid_generate_v4(),
    version        bigserial NOT NULL,
    type_id        int       NOT NULL,
    auth_entity_id uuid      NOT NULL,
    data           jsonb     NOT NULL DEFAULT '{}',
    updated_at     timestamp NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id, version),
    FOREIGN KEY (type_id) REFERENCES item_type (id),
    FOREIGN KEY (auth_entity_id) REFERENCES auth_entity (id)
);
