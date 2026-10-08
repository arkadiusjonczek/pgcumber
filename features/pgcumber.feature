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
