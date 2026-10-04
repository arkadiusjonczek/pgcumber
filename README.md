# pgcumber

![Release Pipeline Status](https://img.shields.io/github/actions/workflow/status/arkadiusjonczek/pgcumber/release.yaml?style=flat-square "Release Pipeline Status")
![Release Version](https://img.shields.io/github/v/release/arkadiusjonczek/pgcumber.svg?style=flat-square&color=blue "Release Version")
![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat-square&color=blue "License")
![Commit Activity](https://img.shields.io/github/commit-activity/m/arkadiusjonczek/pgcumber.svg?style=flat-square&color=blue "Commit Activity")
![Last Commit](https://img.shields.io/github/last-commit/arkadiusjonczek/pgcumber.svg?style=flat-square&color=blue "Last Commit")

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

Use `pgcumber` and add the feature file or directory path as argument:

```shell
pgcumber features/pgcumber.feature
```

### Docker

You can also use the Docker to run `pgcumber` as container:

```shell
docker run \
  --interactive
  --tty
  --rm \
  --volume "$(pwd)/features:/data/pgcumber/features" \
  --network host \
  arkadiusjonczek/pgcumber /data/pgcumber/features
```

## Docker Test Environment

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
