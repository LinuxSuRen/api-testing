*** Settings ***
Library         RequestsLibrary

*** Test Cases ***
one
    ${headers}=    Create Dictionary   key1    value1
    ${response}=    GET  http://foo   headers=${headers}
two
    ${headers}=    Create Dictionary   key2    value2
    ${response}=    GET  http://foo   headers=${headers}
