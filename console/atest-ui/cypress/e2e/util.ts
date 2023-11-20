export function userID_Alpha() {
  let text = "";
  const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

  for (let i = 0; i < 10; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

  return text;
}

export function createTestSuite(store: string, kind: string, name: string, address: string) {
  cy.get('[test-id="open-new-suite-dialog"]').click()

  cy.get('[test-id=suite-form-store] input').click()
  cy.get('[test-id=suite-form-store] input').type(store).trigger('keydown', {key: 'Enter'})

  searchWithID('suite-form-kind', kind)

  cy.get('[test-id=suite-form-name]').should('be.visible').type(name)
  cy.get('[test-id=suite-form-api]').should('be.visible').type(address)
  cy.get('[test-id=suite-form-submit]').should('be.visible').click()
}

export function searchWithID(id: string, txt: string) {
  cy.get('[test-id=' + id + '] input').click()
  cy.get('[test-id=' + id + '] input').type(txt).trigger('keydown', {key: 'Enter'})
}

export function createTestCase(suite: string, name: string, api: string) {
  selectTestSuite(suite)

  cy.get('[test-id="open-new-case-dialog"]').click()
  cy.get('[test-id="case-form-name"]').should('be.visible').type(name)

  cy.get('[test-id="case-form-api"]').should('be.visible').type(api)
  cy.get('[test-id="case-form-submit"]').click()
}

export function selectTestSuite(name: string) {
  cy.get('[test-id="search"]').should('be.visible').type(name)
  cy.get('span').contains(name).should('be.visible').click()
}

export function selectTestCase(name: string) {
  cy.get('[test-id="search"]').should('be.visible').type(name)

  // select the target case
  cy.get('span').contains(name).should('be.visible').click()
}

export function selectTestCaseTab(name: string) {
  cy.get('[id="' + name + '"]').click()
}
