*** Settings ***
Library         RequestsLibrary

*** Test Cases ***
simple
    ${response}=    GET  http://foo
