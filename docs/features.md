# Features of pgcumber

The following list of features are supported by `pgcumber`.

## Autoconfig

Instead of setting each database connection detail (host, port, database, username, password) in every scenario in the features file,
you can use the Autoconfig feature to prevent boilerplate code.

The Autoconfig is always active whenever you set the expected environment variables.

Either create a file and set the file path in the `PGPASSFILE` environment variable with the format:

```text
host:port:database:user:password
``` 

or set the following environment variables:

- `PGHOST`
- `PGPORT`
- `PGUSER`
- `PGPASSWORD`
- `PGDATABASE`

which will overwrite values of the `PGPASSFILE`.

## Autoconnect

The Autoconnect feature is an extension of the Autoconfig feature.

If active it will connect to the database with the given database connection details of the Autoconfig before checking the scenario.

Either use the CLI parameter like:

`pgcumber --autoconnect my.feature`

or use the `@autoconnect` tag like:

```text
@autoconnect
Scenario: This is my scenario
```

To prevent boilerplate it is better to use the `--autoconnect` parameter.
