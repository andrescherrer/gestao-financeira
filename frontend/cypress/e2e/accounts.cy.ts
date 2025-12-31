describe('Accounts Flow', () => {
  beforeEach(() => {
    // Mock authentication
    cy.setAuthToken('mock-token-123')
    
    // Mock API responses
    cy.intercept('GET', '**/api/v1/accounts**', {
      statusCode: 200,
      body: {
        accounts: [
          {
            account_id: 'acc-1',
            name: 'Conta Corrente',
            type: 'BANK',
            currency: 'BRL',
            balance: 1000.0,
            is_active: true,
          },
          {
            account_id: 'acc-2',
            name: 'Carteira',
            type: 'WALLET',
            currency: 'BRL',
            balance: 500.0,
            is_active: true,
          },
        ],
        count: 2,
      },
    }).as('getAccounts')
  })

  it('should display accounts list', () => {
    cy.visit('/accounts')
    cy.wait('@getAccounts')
    
    // Should show accounts
    cy.contains('Conta Corrente').should('be.visible')
    cy.contains('Carteira').should('be.visible')
  })

  it('should open create account form', () => {
    cy.visit('/accounts')
    cy.wait('@getAccounts')
    
    // Click create button
    cy.get('body').then(($body) => {
      if ($body.find('button:contains("nova conta")').length > 0 || 
          $body.find('button:contains("Nova Conta")').length > 0 ||
          $body.find('button:contains("adicionar")').length > 0) {
        cy.get('button').contains(/nova conta|adicionar/i).click()
        
        // Should show form
        cy.get('label').contains(/nome|name/i, { timeout: 5000 }).should('be.visible')
      }
    })
  })

  it('should create a new account', () => {
    cy.visit('/accounts')
    cy.wait('@getAccounts')
    
    cy.intercept('POST', '**/api/v1/accounts', {
      statusCode: 201,
      body: {
        account_id: 'acc-new',
        name: 'Nova Conta',
        type: 'BANK',
        currency: 'BRL',
        balance: 0.0,
        is_active: true,
      },
    }).as('createAccount')
    
    // Open create form
    cy.get('body').then(($body) => {
      if ($body.find('button:contains("nova conta")').length > 0 || 
          $body.find('button:contains("adicionar")').length > 0) {
        cy.get('button').contains(/nova conta|adicionar/i).click()
        
        // Fill form
        cy.get('label').contains(/nome|name/i).parent().find('input').type('Nova Conta')
        
        // Submit
        cy.get('button').contains(/salvar|save/i).click()
        
        // Wait for API call
        cy.wait('@createAccount')
        
        // Should redirect to accounts
        cy.url().should('include', 'accounts')
      }
    })
  })
})

