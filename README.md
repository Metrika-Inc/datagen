# Metrika Datagen
This service consists of a single script that serves a `.json` file in chunks via `/data` endpoint.

## API
* `GET /data` returns an array of JSON encoded data

## Building and running
The service can be built and run as a docker container.

```
make build && make run
```

### Other commands
Check service logs:
```
make logs
```
Stop the service:
```
make stop
```

### Parameters
Service parameters can be adjusted in the makefile:
* `JSON_FILE`: relative path to json file to be served (default: `ledger.json`)
* `HOST_PORT`: host port the Docker service will attach to (default: `9000`)

### Flags
* `-f` relative file path, should match `JSON_FILE`
* `-r` chunks of data to process (default: `149`)