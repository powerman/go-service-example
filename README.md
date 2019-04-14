# Example Go service using go-swagger and Clean Architecture

This project shows an example of how to use go-swagger accordingly to
Uncle Bob's "Clean Architecture".

Also it includes [json-schema-cheatsheet.yml](go-swagger JSON Schema
support cheatsheet), which list all available validations/annotations for
JSON body.

## Clean Architecture

It's not a complete example of Clean Architecture itself (business-logic
of this example is too trivial, so "Use Cases" layer in package `app`
embeds "Entities" layer and implements "DB" instead of have it injected as
dependency through "Gateway"), but it does show relevant part: how to
create "API Controller" layer in package `api` between code auto-generated
by go-swagger and "Use Cases" layer in package `app`.

## Requirements

For using `--strict` swagger option you'll need to install a fork:

```
git clone https://github.com/Djarvur/go-swagger && cd go-swagger && go install
```

## Usage

Command used to (re-)generate go-swagger files:

```
swagger generate server --exclude-main --regenerate-configureapi \
    --target internal/api --api-package op --model-package model \
    --principal app.Auth \
    --strict
```

- We use custom `main.go`.
- It replaces `internal/api/restapi/configure_*.go` because sometimes
  incompatible changes in `swagger.yml` will result in compile error if
  this file won't be updated manually, but since we don't touch it at all
  it's always safe to re-write it.
- It store generated files in `internal/api/restapi/` and
  `internal/api/model/` and you shouldn't need to add or change anything
  in these directories, so they're always safe to remove and re-generate.
- It renames package `operations` to `op` and package `models` to `model`.
- It uses external type `app.Auth` to store authentication details.
- It uses `--strict` from mentioned above fork for better type safety in
  API handlers (each handler is restricted to returning only responses
  defined for that handler instead of general `middleware.Responder`).

## Features

- [X] By default uses host:port from swagger.yml.
- [X] Defaults can be overwritten by env vars and then flags.
- [X] Nice logging with [structlog](https://github.com/powerman/structlog).
- [X] Example go-swagger authentication and authorization.
- [X] CORS, so you can play with it using Swagger Editor tool.
- [X] Easily testable code (as it should be with Clean Architecture).
- [ ] Example tests.
- [ ] Production logging.
- [ ] Production metrics using prometheus.

## Run

Using `./build` script is optional, it's main feature is embedding git
version into compiled binary.

```
$ ./build
$ ./bin/address-book -h
Usage of ./bin/address-book:
  -host host
    	listen on host (default "localhost")
  -log.level level
    	log level (debug|info|warn|err) (default "debug")
  -port port
    	listen on port (>0) (default 8765)
  -version
    	print version
$ ./bin/address-book -version
address-book v0.1.0 4c5dc1b 2019-04-14_17:16:29 go1.12.3
$ ./bin/address-book
address-book[9854] inf   main: `started` version v0.1.0 4c5dc1b 2019-04-14_17:16:29
address-book[9854] inf   main: `protocol` version 0.1.0
address-book[9854] inf   main: `Serving address book at http://127.0.0.1:8765`
```
