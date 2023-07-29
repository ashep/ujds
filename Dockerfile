ARG ARCH
FROM --platform=linux/${ARCH} golang:alpine AS build

ARG APP_NAME
ARG APP_VERSION

WORKDIR /build
COPY . .

RUN mkdir out && \
    go mod vendor && \
    go build -ldflags="-X 'main.buildName=${APP_NAME}' -X 'main.buildVer=${APP_VERSION}'" -o out/app  main.go

ARG ARCH
FROM --platform=linux/${ARCH} alpine:latest

WORKDIR /app
COPY --from=build /build/out/app ./
ENTRYPOINT ["/app/app"]
