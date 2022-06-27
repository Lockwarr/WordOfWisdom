Feature: User requesting quote

  Scenario: Request a quote full scenario
    Given a tcp connection exists with server running
    When I send a "ChallengeRequest" to the server
    Then I receive "ChallengeResponse"
    Then I solve the challenge
    And I send a "QuoteRequest" to the server
    Then I receive "QuoteResponse"
    Then the received Quote is "valid"

#TODO Add more scenarios

  Scenario: Request a quote with invalid challenge
    Given a tcp connection exists with server running
    When I send a "ChallengeRequest" to the server
    Then I receive "ChallengeResponse"
    And I send a "QuoteRequest" to the server
    Then the received Quote is "invalid"

