Feature: PostgreSQL admin user permissions

  Scenario: The user admin is a superuser on the pgcumber database

    Given The host "localhost"
    And   The port "5432"
    And   The database "pgcumber"

    When  I use username "admin"
    And   I use password "admin"
    And   I connect to the server

    Then  the connection is successful
    And   I am a database superuser

    Then  the schema "public" exists
    And   I have "USAGE" privilege on schema "public"
    And   I have "CREATE" privilege on schema "public"

    Then  the table "user" exists in schema "public"
    And   I have "SELECT" privilege on table "public.user"
    And   I have "INSERT" privilege on table "public.user"
    And   I have "UPDATE" privilege on table "public.user"
