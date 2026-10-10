# pgcumber

[![Release Pipeline Status](https://img.shields.io/github/actions/workflow/status/arkadiusjonczek/pgcumber/release.yaml?style=flat-square "Release Pipeline Status")](https://github.com/arkadiusjonczek/pgcumber/actions/workflows/release.yaml)
[![Release Version](https://img.shields.io/github/v/release/arkadiusjonczek/pgcumber.svg?style=flat-square&color=blue "Release Version")](https://github.com/arkadiusjonczek/pgcumber/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat-square&color=blue "License")](https://github.com/arkadiusjonczek/pgcumber/blob/main/LICENSE)
[![Commit Activity](https://img.shields.io/github/commit-activity/m/arkadiusjonczek/pgcumber.svg?style=flat-square&color=blue "Commit Activity")](https://github.com/arkadiusjonczek/pgcumber/commits/main/)
[![Last Commit](https://img.shields.io/github/last-commit/arkadiusjonczek/pgcumber.svg?style=flat-square&color=blue "Last Commit")](https://github.com/arkadiusjonczek/pgcumber/commits/main/)
[![Downloads](https://img.shields.io/github/downloads/arkadiusjonczek/pgcumber/total?style=flat-square&color=blue "Downloads")](https://github.com/arkadiusjonczek/pgcumber/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/arkadiusjonczek/pgcumber?style=flat-square&color=blue "Docker Pulls")](https://hub.docker.com/r/arkadiusjonczek/pgcumber)

Use Gherkin to write plain-text, human-readable test cases for your PostgreSQL data.

```gherkin
Feature: PostgreSQL user permissions

  Scenario: The user admin is a superuser on the pgcumber database

    Given The host "localhost"
    And   The port "5432"
    And   The database "pgcumber"

    When  I use username "admin"
    And   I use password "admin"
    And   I connect to the server

    Then  the connection is successful
    And   I am a database superuser

  Scenario: The user johndoe is not a superuser on the pgcumber database

    Given The host "localhost"
    And   The port "5432"
    And   The database "pgcumber"

    When  I use username "johndoe"
    And   I use password "johndoe"
    And   I connect to the server

    Then  the connection is successful
    And   I am not a database superuser
```

## Usage

Use `pgcumber` and add the feature file or directory path as argument:

```shell
pgcumber features/pgcumber.feature
```

### macOS

If you are on macOS and you download the binary in e.g. the web browser you might remove the assigned quarantine attribute:

```shell
xattr -d com.apple.quarantine pgcumber
```

## Docker

You can also use the Docker image from [Docker Hub](https://hub.docker.com/r/arkadiusjonczek/pgcumber) to run `pgcumber` as container:

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
