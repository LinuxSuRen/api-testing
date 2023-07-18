// https://on.cypress.io/api

describe('Components exist', () => {
  it('Home page', () => {
    cy.visit('/')
    cy.contains('span', 'Tool Box')
  })
})
