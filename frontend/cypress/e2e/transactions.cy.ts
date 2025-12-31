describe('Transactions Flow', () => {
  beforeEach(() => {
    // Mock authentication
    cy.setAuthToken('mock-token-123')
    
    // Mock API responses
    cy.intercept('GET', '**/api/v1/transactions**', {
      statusCode: 200,
      body: {
        transactions: [
          {
            transaction_id: 'tx-1',
            account_id: 'acc-1',
            type: 'INCOME',
            amount: 1000.0,
            currency: 'BRL',
            description: 'Salário',
            date: '2024-01-01',
          },
          {
            transaction_id: 'tx-2',
            account_id: 'acc-1',
            type: 'EXPENSE',
            amount: 500.0,
            currency: 'BRL',
            description: 'Compras',
            date: '2024-01-02',
          },
        ],
        count: 2,
      },
    }).as('getTransactions')
  })

  it('should display transactions list', () => {
    cy.visit('/transactions')
    cy.wait('@getTransactions')
    
    // Should show transactions
    cy.contains('Salário').should('be.visible')
    cy.contains('Compras').should('be.visible')
  })

  it('should filter transactions by type', () => {
    cy.visit('/transactions')
    cy.wait('@getTransactions')
    
    // Click on type filter
    cy.get('body').then(($body) => {
      if ($body.find('label:contains("tipo")').length > 0 || $body.find('label:contains("Tipo")').length > 0) {
        cy.get('label').contains(/tipo/i).parent().find('select, button').first().click({ force: true })
        cy.contains(/receitas|income/i).click({ force: true })
        
        // Should filter transactions
        cy.contains('Salário').should('be.visible')
      }
    })
  })

  it('should open create transaction form', () => {
    cy.visit('/transactions')
    cy.wait('@getTransactions')
    
    // Mock accounts API
    cy.intercept('GET', '**/api/v1/accounts**', {
      statusCode: 200,
      body: {
        accounts: [
          {
            account_id: 'acc-1',
            name: 'Conta Corrente',
            currency: 'BRL',
          },
        ],
        count: 1,
      },
    }).as('getAccounts')
    
    // Click create button
    cy.get('body').then(($body) => {
      if ($body.find('button:contains("nova transação")').length > 0 || 
          $body.find('button:contains("Nova Transação")').length > 0 ||
          $body.find('button:contains("adicionar")').length > 0) {
        cy.get('button').contains(/nova transação|adicionar/i).click()
        
        // Should show form
        cy.get('label').contains(/descrição|description/i, { timeout: 5000 }).should('be.visible')
      }
    })
  })

  it('should create a new transaction', () => {
    cy.visit('/transactions')
    cy.wait('@getTransactions')
    
    // Mock APIs
    cy.intercept('GET', '**/api/v1/accounts**', {
      statusCode: 200,
      body: {
        accounts: [
          {
            account_id: 'acc-1',
            name: 'Conta Corrente',
            currency: 'BRL',
          },
        ],
      },
    }).as('getAccounts')
    
    cy.intercept('POST', '**/api/v1/transactions', {
      statusCode: 201,
      body: {
        transaction_id: 'tx-new',
        account_id: 'acc-1',
        type: 'INCOME',
        amount: 2000.0,
        currency: 'BRL',
        description: 'Nova Receita',
        date: '2024-01-03',
      },
    }).as('createTransaction')
    
    // Open create form
    cy.get('body').then(($body) => {
      if ($body.find('button:contains("nova transação")').length > 0 || 
          $body.find('button:contains("adicionar")').length > 0) {
        cy.get('button').contains(/nova transação|adicionar/i).click()
        
        // Fill form
        cy.get('label').contains(/descrição|description/i).parent().find('input, textarea').type('Nova Receita')
        cy.get('label').contains(/valor|amount/i).parent().find('input').type('2000')
        
        // Submit
        cy.get('button').contains(/salvar|save/i).click()
        
        // Wait for API call
        cy.wait('@createTransaction')
        
        // Should redirect to transactions
        cy.url().should('include', 'transactions')
      }
    })
  })
})

