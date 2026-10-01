# pgcumber

Use Gherkin to write plain-text, human-readable test cases for your PostgreSQL data. 

```gherkin
Feature: PostgreSQL

  Scenario: PostgreSQL connection & user

    Given The host "localhost"
    And   The port "5432"
    And   The database "pgcumber"

    When I use username "pgcumberuser"
    And  I use password "pgcumberpassword"
    And  I connect to the server

    Then the connection is successful
    And  I am a database superuser
```

## Usage

Use the cli tool `pgcumber` and add the feature file path oder the feature files directory path as argument:

```shell
pgcumber features/pgcumber.feature
```

## Build

Use go to build the `pgcumber` cli tool:

```shell
go build -o pgcumber ./cmd/pgcumber
```

## Test

Use the docker compose file [docker-compose.yaml](docker-compose.yaml) to start PostgreSQL:

```shell
docker-compose up -d
```

Then use `pgcumber` and the feature file [features/pgcumber.feature](features/pgcumber.feature):

```shell
pgcumber features/pgcumber.feature
```

After your tests you can stop the docker environment again:

```shell
docker-compose down
```

## Background

The implementation is based on [godog](https://github.com/cucumber/godog) which is cucumber for golang.