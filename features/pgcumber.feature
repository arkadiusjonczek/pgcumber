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
