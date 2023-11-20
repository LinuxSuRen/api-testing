import { createTestCase, createTestSuite, selectTestCase, selectTestCaseTab, userID_Alpha } from './util'

describe('gRPC', () => {
  const suiteName = userID_Alpha()
  const store = "local"
  const sampleAPIAddress = "http://foo"

  it('Create Suite', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()

    createTestSuite(store, 'gRPC', suiteName, sampleAPIAddress)
  })

  const caseName = userID_Alpha()
  const caseAPI = "/api/v2"
  it('New Test Case', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()
    createTestCase(suiteName, caseName, caseAPI)
  })

  it('Update Test Case', () => {
    cy.visit('/')
    cy.get('.introjs-skipbutton').click()
    selectTestCase(caseName)
    selectTestCaseTab('tab-second')
  })
})
