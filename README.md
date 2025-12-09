# UJDS: a Universal JSON Data Storage

[![ci](https://github.com/ashep/ujds/actions/workflows/ci.yaml/badge.svg)](https://github.com/ashep/ujds/actions/workflows/ci.yaml)

The **Universal JSON Data Storage** stores arbitrary JSON data and keeps changes history. Data are being stored in
**indices** as **records**. Indices may have **schema** to check incoming data upon updates.

## Configuration

The service can be configured in three ways: via YAML file, via env variables or using both. Env variables take
precedence over config file.

If the `config.yaml` file is found in the current directory, it will be loaded before env variables. It is possible to
change default config file location using `APP_CONFIG_PATH` env variable.

### File

- *required* **object** `db`: database configuration.
    - *required* **string** `dsn`: database source name.
- *optional* **object** `server`: server configuration.
    - *optional* **string** `address`: network address, default is `:9000`.
    - *optional* **string** `auth_token`: authorization token.

### Env variables

- *required* **string** `UJDS_DB_DSN`: database source name.
- *optional* **string** `UJDS_SERVER_ADDRESS`: server network address.
- *optional* **string** `UJDS_SERVER_AUTHTOKEN`: server authorization token.

## HTTP API

### Methods

All the requests must be performed use `POST` method.

### Response JSON types

Please note that **numerical data in responses are encoded as strings**.

### Authorization

If the `server.auth_token` configuration parameter is specified, the server will expect an `Authorization: Bearer XXX`
HTTP header, where `XXX` must match the configured token value, otherwise the `403` HTTP status will be returned.

### Error handling

If a request is not successful, the service responds with an HTTP status code other than 200, providing a JSON object
with the `code` and `message` fields. Use that information to understand what went wrong. Example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.index.v1.IndexService/Push \
  --header 'Authorization: Bearer WrongAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{}'
```

```text
< HTTP/2 401
< content-type: application/json
<
{"code":"unauthenticated","message":"not authorized"}
```

### Search query syntax

The `RecordService/Find` method provides a method of filtering result using search queries. The syntax has to be
described here.

### IndexService/Push

Creates a new index or updates an existing one.

- Request fields:
    - *required* **string** `name`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.
    - *optional* **string** `title`: index title.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.index.v1.IndexService/Push \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "books",
	"title": "The books"
}'
```

### IndexService/Get

Returns an index metadata.

- Request fields:
    - *required* **string** `name`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.
- Response fields:
    - **string** `name`: index name.
    - **string** `title`: index title.
    - **int** `createdAt`: creation UNIX timestamp.
    - **int** `updatedAt`: update UNIX timestamp.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.index.v1.IndexService/Get \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{"name": "books"}'
```

Response example:

```json
{
  "name": "books",
  "title": "The books",
  "createdAt": "1693768684",
  "updatedAt": "1693769057"
}
```

### IndexService/List

Returns existing indices list.

- Request fields:
    - *optional* **object** `filter`: filter.
        - *optional* **[]string** `names`: index name patterns. Allowed wildcard symbols: `*`.
- Response fields:
    - **[]object** `indices`
        - **string** `name`: index name.
        - **string** `title`: index title.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.index.v1.IndexService/List \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{"filter":{"names": ["book*", "recip*", "cartoons"]}}'
```

Response example:

```json
{
  "indices": [
    {
      "name": "books",
      "title": "The books"
    },
    {
      "name": "recipes",
      "title": "The recipes"
    },
    {
      "name": "cartoons",
      "title": "The cartoons"
    }
  ]
}
```

### IndexService/Clear

Clears all index records.

- Request fields:
    - *required* **string** `name`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.index.v1.IndexService/Clear \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{}'
```

### RecordService/Push

Creates records in the index or updates existing ones.

- Request fields:
    - *required* **[]object** `records`: records.
        - *required* **string** `index`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.
        - *required* **string** `id`: record ID.
        - *required* **string** `data`: record JSON data.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.record.v1.RecordService/Push \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{
	"records": [
		{
			"index": "books",
			"id": "castaneda-001",
			"data": "{\"title\": \"Tales of Power\", \"author\": \"Carlos Castaneda\", \"isbn\":\"978-0-671-73252-3\"}"
		},
		{
			"index": "books",
			"id": "tanenbaum-001",
			"data": "{\"author\":\"M. van Steen and A.S. Tanenbaum\", \"title\":\"Distributed Systems, 4th ed.\"}"
		}
	]
}'
```

### RecordService/Get

Returns a single record.

- Request fields:
    - *required* **string** `index`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.
- Response field:
    - **object** `record`
        - **string** `id`: ID.
        - **string** `index`: index name.
        - **string** `rev`: revision number.
        - **string** `createdAt`: creation time as UNIX timestamp.
        - **string** `updatedAt`: last change time as UNIX timestamp.
        - **string** `touchedAt`: last update time as UNIX timestamp.
        - **string** `data`: data.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.record.v1.RecordService/Get \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{
	"index": "books",
	"id": "castaneda-001"
}'
```

Response example:

```json
{
  "record": {
    "id": "castaneda-001",
    "rev": "227",
    "index": "books",
    "createdAt": "1694109017",
    "updatedAt": "1694237265",
    "touchedAt": "1702938162",
    "data": "{\"title\": \"Tales of Power\", \"author\": \"Carlos Castaneda\", \"isbn\":\"978-0-671-73252-3\"}"
  }
}
```

### RecordService/Find

Returns all records from the index.

- Request fields:
    - *required* **string** `index`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.
    - *optional* **string** `search`: search query. TODO: describe search query syntax.
    - *optional* **int** `since`: return only records, that have been **modified** since provided UNIX timestamp.
    - *optional* **int** `touchedSince`: return only records, that have been **touched** since a UNIX timestamp.
    - *optional* **int** `notTouchedSince`: return only records, that have not been **touched** since a UNIX timestamp.
    - *optional* **int** `cursor`: pagination: return records starting from provided position.
    - *optional* **int** `limit`: get only specified number of records; default and maximum is `500`.
- Response fields:
    - **string** `cursor`: pagination cursor position, that should be used to retrieve the next result set.
    - **[]object** `records`
        - **string** `id`: ID.
        - **string** `index`: index name.
        - **string** `rev`: revision number.
        - **string** `createdAt`: creation time as UNIX timestamp.
        - **string** `updatedAt`: last change time as UNIX timestamp.
        - **string** `touchedAt`: last update time as UNIX timestamp.
        - **string** `data`: data.

Request example:

```shell
curl --request POST \
  --url https://localhost:9000/ujds.record.v1.RecordService/Find \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{
	"index": "books",
	"search": "author=\"Carlos Castaneda\"",
	"since": 1694109017,
	"cursor": 226,
	"limit": 2
}'
```

Response example:

```json
{
  "cursor": "228",
  "records": [
    {
      "id": "castaneda-001",
      "rev": "227",
      "index": "books",
      "createdAt": "1694109017",
      "updatedAt": "1694109017",
      "touchedAt": "1702938162",
      "data": "{\"title\": \"Tales of Power\", \"author\": \"Carlos Castaneda\", \"isbn\":\"978-0-671-73252-3\"}"
    },
    {
      "id": "castaneda-002",
      "rev": "228",
      "index": "books",
      "createdAt": "1694109017",
      "updatedAt": "1694109017",
      "touchedAt": "1702938162",
      "data": "{\"title\": \"The Fire From Within\", \"author\": \"Carlos Castaneda\", \"isbn\":\"978-0-671-73250-9\"}"
    }
  ]
}
```

### RecordService/History

Returns record history.

- Request fields:
    - *required* **string** `index`: index name. The allowed format: `^[a-zA-Z0-9.-]{1,255}$`.
    - *required* **string** `id`: record id.
    - *optional* **int** `since`: return only history records, which have been created since provided UNIX timestamp.
    - *optional* **int** `cursor`: pagination: return records starting from provided position.
    - *optional* **int** `limit`: get only specified number of records; default and maximum is `500`.
- Response fields:
    - **string** `cursor`: pagination cursor position, which should be used to retrieve the next result set.
    - **[]object** `records`
        - **string** `id`: record ID.
        - **string** `index`: index name.
        - **string** `rev`: revision number.
        - **string** `createdAt`: creation time as UNIX timestamp.
        - **string** `data`: data.

Request example:

```shell
curl --request POST \
  --url http://localhost:9000/ujds.record.v1.RecordService/History \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{
	"index": "books",
	"id": "castaneda-001",
	"since": 1696767680,
	"cursor": 28,
	"limit": 2
}'
````

Response example:

```json
{
  "records": [
    {
      "id": "castaneda-001",
      "rev": "30",
      "index": "books",
      "createdAt": "1696768530",
      "data": "{\"title\": \"Tales of Power, second edition\", \"author\": \"Carlos Castaneda\", \"isbn\":\"978-0-671-73252-3\"}"
    },
    {
      "id": "castaneda-001",
      "rev": "28",
      "index": "books",
      "createdAt": "1696767687",
      "data": "{\"title\": \"Tales of Power\", \"author\": \"Carlos Castaneda\", \"isbn\":\"978-0-671-73252-3\"}"
    }
  ]
}
```

## Developers notes

Create migration:

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext .sql -dir internal/migration/migrations foobar
```

## Changelog

### 0.9.2 (2025-12-09)

Fixed config parsing.

### 0.9.1 (2025-12-09)

Fixed tests.

### 0.9 (2025-12-09)

Add config-based per-index schema validation.

### 0.8 (2025-12-05)

- Get rid of per-index schema.
- Global schema introduced.

### 0.7.1 (2025-12-04)

Fixed HTTP server default address configuration.

### 0.7 (2025-12-04)

`RecordService/Find` request got the new `touchedSince` argument.


### 0.6 (2025-09-02)

`RecordService/Find` request got the new `notTouchedSince` argument.


### 0.5 (2024-03-30)

The `filter` request field added to `IndexService/List` method.

### 0.4 (2024-03-09)

- `RecordService/Push`:
    - now the index should be specified on each record;
    - a new `touched_at` field added; it is always being updated with a current timestamp, even if record hasn't
      changed.

### 0.3 (2023-11-29)

- `RecordService/GetAll` API method renamed to `RecordService/Find`.
- A new `search` optional field added to `RecordService/Find`.

### 0.2 (2023-10-08)

- `RecordService/History` API method added.
- Index name length extended to 255 chars.
- Slash is not allowed in index names anymore; replaced with dot.

### 0.1 (2023-09-07)

Initial release.

## Authors

- [Oleksandr Shepetko](https://shepetko.com)
