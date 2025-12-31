import 'cypress-axe'

describe('Accessibility Tests', () => {
  beforeEach(() => {
    // Mock authentication
    cy.setAuthToken('mock-token-123')
  })

  it('should have no accessibility violations on login page', () => {
    cy.visit('/login')
    cy.injectAxe()
    cy.checkA11y()
  })

  it('should have no accessibility violations on dashboard', () => {
    // Mock API responses
    cy.intercept('GET', '**/api/v1/**', {
      statusCode: 200,
      body: {},
    })
    
    cy.visit('/')
    cy.injectAxe()
    cy.checkA11y()
  })

  it('should have no accessibility violations on transactions page', () => {
    // Mock API responses
    cy.intercept('GET', '**/api/v1/transactions**', {
      statusCode: 200,
      body: {
        transactions: [],
        count: 0,
      },
    })
    
    cy.visit('/transactions')
    cy.injectAxe()
    cy.checkA11y()
  })

  it('should have proper heading hierarchy', () => {
    cy.visit('/login')
    
    // Check for h1
    cy.get('h1').should('have.length.at.least', 1)
  })

  it('should have proper form labels', () => {
    cy.visit('/login')
    
    // Check for form labels
    cy.get('label').contains(/email/i).should('be.visible')
    cy.get('label').contains(/password|senha/i).should('be.visible')
  })

  it('should have proper button labels', () => {
    cy.visit('/login')
    
    // Check for button with accessible name
    cy.get('button').contains(/entrar|login/i).should('be.visible')
  })
})

