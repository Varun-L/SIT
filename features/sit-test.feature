Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 3 out of 12
    Given there are 12 godogs
    When I eat 3
    Then there should be 9 remaining
  
    Scenario: Eat 2 out of 12
    Given there are 12 godogs
    When I eat 2
    Then there should be 10 remaining
  
    Scenario: Eat 2 out of 12
    Given there are 12 godogs
    When I eat 10
    Then there should be 0 remaining