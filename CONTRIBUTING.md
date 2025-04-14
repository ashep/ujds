# DataPimp Developer's Documentation

## Local run

```shell
task run
```

or

```shell
go run -race internal/main/main.go
```

## Generate protobuf code

Make sure you have installed necessary tools:

```shell
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

Then run:

```shell
task gen
```

or

```shell
buf lint && buf generate
```

## Create a new migration scripts

Make sure you have installed necessary tools:

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Create and edit new migration scripts:

```shell
migrate create -dir migration/migrations -ext sql description
```

Edit the new scripts.
