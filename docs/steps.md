# Available Steps for pgcumber

You can use the following Steps as described in Given, When or Then Steps.

## Given

```
Given The host "localhost"
```

```
Given The port "5432"
```

```
Given The database "pgcumber"
```

## When

```
When I use username "admin"
```

```
When I use username from env "PSQL_USERNAME"
When I use username from env var "PSQL_USERNAME"
When I use username from environment variable "PSQL_USERNAME"
```

```
When I use username from file "psql_username"
```

```
When I use password "admin"
```

```
When I use password from env "PSQL_PASSWORD"
When I use password from env var "PSQL_PASSWORD"
When I use password from environment variable "PSQL_PASSWORD"
```

```
When I use password from file "psql_password"
```

```
When I connect to the server
```

## Then

```
Then the connection is successful
```

```
Then I am a database superuser
```

```
Then I am not a database superuser
```
