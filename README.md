# UJDS

[![ci](https://github.com/ashep/ujds/actions/workflows/ci.yaml/badge.svg)](https://github.com/ashep/ujds/actions/workflows/ci.yaml)

The **Universal JSON Data Storage** is a service aimed to store arbitrary JSON data keeping changes history. Data are
being stored in **indices** as **records**. Indices may be supplied with JSON **schemas** to check incoming data upon 
updates.

## Configuration

The service can be configured in three ways: via YAML file, via env variables or using both. Env variables takes
precedence over config file.

If the `config.yaml` file is found in the current directory, it will be loaded before env variables. It is possible to
change default config file location using `APP_CONFIG_PATH` env variable.

### File

- *required* **object** `db`: database configuration:
    - *required* **string** `dsn`: database source name.
- *optional* **object** `server`: server configuration:
    - *optional* **string** `address`: network address, default is `:9000`;
    - *optional* **string** `auth_token`: authorization token.

### Env variables

- *required* **string** `UJDS_DB_DSN`: database source name.
- *optional* **string** `UJDS_SERVER_ADDRESS`: server network address.
- *optional* **string** `UJDS_SERVER_AUTHTOKEN`: server authorization token.

## HTTP API

### Methods

All the requests must be performed use `POST` method.

### Authorization

If the `server.auth_token` is specified, the server will expect an `Authorization: Bearer XXX` HTTP header, where `XXX`
must match the configured token value, otherwise the `403` HTTP status will be returned.

### Error handling

If a request is not successful, the service responds with an HTTP status code other than 200, providing a JSON object 
with the `code` and `message` fields. Use that information to understand what went wrong. Example:

```shell
curl --request POST \
  --url https://test.com/ujds.index.v1.IndexService/Push \
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

### Push index

Create a new index or update an existing one.

- Path: `/ujds.index.v1.IndexService/Push`
- Request fields:
    - *required* **string** `name`: index name. The allowed format: `^[a-zA-Z0-9_-]{1,64}$`.
    - *optional* **string** `schema`: JSON validation schema.

Request example:

```shell
curl --request POST \
  --url https://test.com/ujds.index.v1.IndexService/Push \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "books",
	"schema": "{\"required\":[\"author\",\"title\"]}"
}'
```

### Get index

- Path: `/ujds.index.v1.IndexService/Get`
- Request fields:
    - *required* **string** `name`: index name. The allowed format: `^[a-zA-Z0-9_-]{1,64}$`.
- Response fields:
  - **string** `name`: index name;
  - **int** `createdAt`: creation UNIX timestamp;
  - **int** `updatedAt`: update UNIX timestamp;
  - **string** `schema`: JSON validation schema.

Request example:

```shell
curl --request POST \
  --url https://test.com/ujds.index.v1.IndexService/Get \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{"name": "books"}'
```

Response example:

```json
{
  "name": "books",
  "createdAt": "1693768684",
  "updatedAt": "1693769057",
  "schema": "{\"required\": [\"author\", \"title\"]}"
}
```

## Get indices list

- Path: `/ujds.index.v1.IndexService/List`
- Response fields:
  - **[]string** `indices`: index names.


Request example:

```shell
curl --request POST \
  --url https://test.com/ujds.index.v1.IndexService/List \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{}'
```

Response example:

```json
{
  "indices": [
    {
      "name": "books"
    },
    {
      "name": "foo"
    },
    {
      "name": "bar"
    }
  ]
}
```

## Clear index

- Path: `/ujds.index.v1.IndexService/Clear`
- Request fields:
  - *required* **string** `name`: index name. The allowed format: `^[a-zA-Z0-9_-]{1,64}$`.

Request example:

```shell
curl --request POST \
  --url https://test.com/ujds.index.v1.IndexService/Clear \
  --header 'Authorization: Bearer YourAuthToken' \
  --header 'Content-Type: application/json' \
  --data '{}'
```

## Changelog

## Authors

- [Oleksandr Shepetko](https://shepetko.com)
