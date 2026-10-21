Feature: PostgreSQL johndoe user permissions

  Scenario: The user johndoe is not a superuser on the pgcumber database

    Given The host "localhost"
    And   The port "5432"
    And   The database "pgcumber"

    When  I use username "johndoe"
    And   I use password "johndoe"
    And   I connect to the server

    Then  the connection is successful
    And   I am not a database superuser
