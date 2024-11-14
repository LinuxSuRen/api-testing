*** Settings ***
Library         RequestsLibrary

*** Test Cases ***
simple
    ${headers}=    Create Dictionary   key    value
    ${response}=    GET  http://foo   headers=${headers}
