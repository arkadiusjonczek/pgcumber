# Available Steps for pgcumber

You can use the following Steps separated in Given, When and Then Steps.

## Given

### Given The host "\$1"

Example:
```
Given The host "localhost"
```

### Given The port "\$1" 

Example:
```
Given The port "5432"
```

### Given The database "\$1"

Example:
```
Given The database "pgcumber"
```

## When

### When I use username "\$1"

Example:
```
When I use username "admin"
```

### When I use username from environment variable "\$1"

Example:
```
When I use username from environment variable "PSQL_USERNAME"
```

### When I use username from file "\$1"

Example:
```
When I use username from file "psql_username"
```

### When I use password "\$1"

Example:
```
When I use password "admin"
```

### When I use password from environment variable "\$1"

Example:
```
When I use password from environment variable "PSQL_PASSWORD"
```

### When I use password from file "\$1"

Example:
```
When I use password from file "psql_password"
```

### When I connect to the server

Example:
```
When I connect to the server
```

## Then

### Then the connection is successful

Example:
```
Then the connection is successful
```

### Then I am a database superuser

Example:
```
Then I am a database superuser
```

### Then I am not a database superuser

Example:
```
Then I am not a database superuser
```

### Then the schema "\$1" exists

Example:
```
Then the schema "public" exists
```

### Then the user has "\$1" privilege on schema "\$2"

Example:
```
Then the user has "USAGE" privilege on schema "public"
```

### Then the table "\$1" exists in schema "\$2"

Example:
```
Then the table "user" exists in schema "public"
```

### Then the current user has "\$1" privilege on table "\$2"

Example:
```
Then the current user has "SELECT" privilege on table "public.user"
```
