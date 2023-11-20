// https://on.cypress.io/api
// https://docs.cypress.io/api/table-of-contents
import { createTestCase, createTestSuite, selectTestCase, selectTestSuite, userID_Alpha } from './util'

describe('Suite Manage', () => {
  const suiteName = userID_Alpha()
  const store = "local"
  const sampleAPIAddress = "http://foo"

  it('Create Suite', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()

    cy.contains('span', 'Tool Box')
    createTestSuite(store, 'HTTP', suiteName, sampleAPIAddress)
  })

  const caseName = userID_Alpha()
  const caseMethod = "GET"
  const caseAPI = "/api/v2"
  it('New Test Case', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()
    createTestCase(suiteName, caseName, caseAPI)
  })

  it('Find Test Case', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()
    selectTestCase(caseName)

    // verify the value
    cy.get('[test-id="case-editor-method"] input').should('have.value', caseMethod)
  })

  it('Delete Suite', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()
    selectTestSuite(suiteName)

    cy.get('[test-id="suite-editor-api"]').should('have.value', sampleAPIAddress)
    cy.get('[test-id="suite-del-but"]').click()
  })
})
