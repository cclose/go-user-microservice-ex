## Prerequisites

You must have docker installed, compatible with docker compose v3.7

## Set up

Create a `.env` file next to `docker-compose.yml`. Replace `<username>` and `<password>` with your desired values:
```
PG_USER_USER=<username>
PG_USER_PASSWORD=<password>
```

This file is read by the Dockerfile to distribute credentials into the containers safely

## Build

Build the App stack with:

    docker-composer up --build

## Tests 

### Unit Test
Unit tests for the User model's validation functions can be found in the `models` folder. With go installed run:

    cd user-service/src/models; go test

If attempting to run the unit test within the docker container, you must first install gcc with:

    apk add --no-cache gcc musl-dev

gcc is required for `go test` to run, but is left out of the base image to keep it lightweight. It is best to run `go test` from your development workspace

### Integration Test

Visit http://localhost:8080/test to see a Javascript Integration Test.