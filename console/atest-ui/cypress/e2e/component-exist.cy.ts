// https://on.cypress.io/api
// https://docs.cypress.io/api/table-of-contents

console.log(Cypress.browser)
describe('Suite Manage', () => {
  const suiteName = userID_Alpha()
  const sampleAPIAddress = "http://foo"
  console.log(sampleAPIAddress)
  
  it('Create Suite', () => {
    cy.visit('/')
    cy.contains('span', 'Tool Box')

    cy.get('[test-id="open-new-suite-dialog"]').click()

    cy.get('[test-id=suite-form-name]').should('be.visible').type(suiteName)
    cy.get('[test-id=suite-form-api]').should('be.visible').type(sampleAPIAddress)
    cy.get('[test-id=suite-form-submit]').should('be.visible').click()
  })

  const caseName = userID_Alpha()
  const caseMethod = "POST"
  const caseAPI = "/api/v2"
  it('New Case', () => {
    cy.visit('/')
    cy.get('span').contains(suiteName).should('be.visible').click()

    cy.get('[test-id="open-new-case-dialog"]').click()
    cy.get('[test-id="case-form-name"]').should('be.visible').type(caseName)
    const methodSelector = cy.get('[test-id="case-form-method"]')
    methodSelector.clear()
    methodSelector.type(caseMethod)

    cy.get('[test-id="case-form-api"]').should('be.visible').type(caseAPI)
    cy.get('[test-id="case-form-submit"]').click()
  })

  it('Find Case', () => {
    cy.visit('/')
    const searchInput = cy.get('[test-id="search"]')
    searchInput.type(caseName)
    searchInput.trigger('keydown', {key: 'Enter'})

    // select the target case
    cy.get('span').contains(caseName).should('be.visible').click()

    // verify the value
    cy.get('[test-id="case-editor-method"] input').should('have.value', caseMethod)
  })

  it('Delete Suite', () => {
    cy.visit('/')
    cy.get('span').contains(suiteName).should('be.visible').click()
    cy.get('[test-id="suite-editor-api"]').should('have.value', sampleAPIAddress)
    cy.get('[test-id="suite-del-but"]').click()
  })
})

function userID_Alpha() {
  let text = "";
  const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

  for (let i = 0; i < 10; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

  return text;
}