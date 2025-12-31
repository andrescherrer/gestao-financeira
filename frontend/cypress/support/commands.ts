/// <reference types="cypress" />

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command to login with mock credentials
       * @example cy.login()
       */
      login(): Chainable<void>
      
      /**
       * Custom command to set auth token in localStorage
       * @example cy.setAuthToken('mock-token-123')
       */
      setAuthToken(token: string): Chainable<void>
      
      /**
       * Custom command to mock API response
       * @example cy.mockApi('GET', '/api/v1/accounts', { accounts: [] })
       */
      mockApi(method: string, url: string, response: any, status?: number): Chainable<void>
    }
  }
}

Cypress.Commands.add('login', () => {
  cy.setAuthToken('mock-token-123')
  cy.window().then((win) => {
    win.localStorage.setItem('auth_token', 'mock-token-123')
  })
})

Cypress.Commands.add('setAuthToken', (token: string) => {
  cy.window().then((win) => {
    win.localStorage.setItem('auth_token', token)
  })
})

Cypress.Commands.add('mockApi', (method: string, url: string, response: any, status = 200) => {
  cy.intercept(method as Cypress.HttpMethod, url, {
    statusCode: status,
    body: response,
  }).as(`mock${method}${url.replace(/\//g, '_')}`)
})

export {}

