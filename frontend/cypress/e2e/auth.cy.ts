describe('Authentication Flow', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('should redirect to login when not authenticated', () => {
    // Try to access a protected route
    cy.visit('/accounts')
    
    // Should redirect to login
    cy.url().should('include', '/login')
  })

  it('should show login form', () => {
    cy.visit('/login')
    
    // Check for login form elements
    cy.get('label').contains(/email/i).should('be.visible')
    cy.get('label').contains(/password|senha/i).should('be.visible')
    cy.get('button').contains(/entrar|login/i).should('be.visible')
  })

  it('should display validation errors on invalid login', () => {
    cy.visit('/login')
    
    // Try to submit empty form
    cy.get('button').contains(/entrar|login/i).click()
    
    // Should show validation errors
    cy.contains(/email.*required|obrigatório/i, { timeout: 5000 }).should('be.visible')
  })

  it('should login successfully with valid credentials', () => {
    cy.visit('/login')
    
    // Mock successful login
    cy.intercept('POST', '**/api/v1/auth/login', {
      statusCode: 200,
      body: {
        token: 'mock-token-123',
        user: {
          user_id: 'user-123',
          email: 'test@example.com',
          first_name: 'Test',
          last_name: 'User',
        },
      },
    }).as('loginRequest')
    
    // Fill login form
    cy.get('label').contains(/email/i).parent().find('input').type('test@example.com')
    cy.get('label').contains(/password|senha/i).parent().find('input').type('password123')
    
    // Submit form
    cy.get('button').contains(/entrar|login/i).click()
    
    // Wait for API call
    cy.wait('@loginRequest')
    
    // Should redirect to dashboard or home
    cy.url().should('match', /.*\/$|.*dashboard/)
  })

  it('should logout successfully', () => {
    // First login
    cy.visit('/login')
    
    // Mock login
    cy.intercept('POST', '**/api/v1/auth/login', {
      statusCode: 200,
      body: {
        token: 'mock-token-123',
        user: {
          user_id: 'user-123',
          email: 'test@example.com',
        },
      },
    }).as('loginRequest')
    
    cy.get('label').contains(/email/i).parent().find('input').type('test@example.com')
    cy.get('label').contains(/password|senha/i).parent().find('input').type('password123')
    cy.get('button').contains(/entrar|login/i).click()
    
    // Wait for navigation
    cy.wait('@loginRequest')
    cy.url().should('match', /.*\/$|.*dashboard/)
    
    // Find and click logout button
    cy.get('body').then(($body) => {
      if ($body.find('button:contains("logout")').length > 0 || $body.find('button:contains("sair")').length > 0) {
        cy.get('button').contains(/logout|sair/i).click()
        
        // Should redirect to login
        cy.url().should('include', '/login')
      }
    })
  })
})

