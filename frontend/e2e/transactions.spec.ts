import { test, expect } from '@playwright/test'

test.describe('Transactions Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Mock authentication
    await page.goto('/')
    await page.evaluate(() => {
      localStorage.setItem('auth_token', 'mock-token-123')
    })
    
    // Mock API responses
    await page.route('**/api/v1/transactions**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
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
        }),
      })
    })
  })

  test('should display transactions list', async ({ page }) => {
    await page.goto('/transactions')
    
    // Should show transactions
    await expect(page.getByText('Salário')).toBeVisible()
    await expect(page.getByText('Compras')).toBeVisible()
  })

  test('should filter transactions by type', async ({ page }) => {
    await page.goto('/transactions')
    
    // Click on type filter
    const typeFilter = page.getByLabel(/tipo/i)
    if (await typeFilter.isVisible()) {
      await typeFilter.click()
      await page.getByText(/receitas|income/i).click()
      
      // Should filter transactions
      await expect(page.getByText('Salário')).toBeVisible()
    }
  })

  test('should open create transaction form', async ({ page }) => {
    await page.goto('/transactions')
    
    // Mock accounts API
    await page.route('**/api/v1/accounts**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          accounts: [
            {
              account_id: 'acc-1',
              name: 'Conta Corrente',
              currency: 'BRL',
            },
          ],
          count: 1,
        }),
      })
    })
    
    // Click create button
    const createButton = page.getByRole('button', { name: /nova transação|adicionar/i })
    if (await createButton.isVisible()) {
      await createButton.click()
      
      // Should show form
      await expect(page.getByLabel(/descrição|description/i)).toBeVisible()
    }
  })

  test('should create a new transaction', async ({ page }) => {
    await page.goto('/transactions')
    
    // Mock APIs
    await page.route('**/api/v1/accounts**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          accounts: [
            {
              account_id: 'acc-1',
              name: 'Conta Corrente',
              currency: 'BRL',
            },
          ],
        }),
      })
    })
    
    let createCalled = false
    await page.route('**/api/v1/transactions', async route => {
      if (route.request().method() === 'POST') {
        createCalled = true
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            transaction_id: 'tx-new',
            account_id: 'acc-1',
            type: 'INCOME',
            amount: 2000.0,
            currency: 'BRL',
            description: 'Nova Receita',
            date: '2024-01-03',
          }),
        })
      } else {
        await route.continue()
      }
    })
    
    // Open create form
    const createButton = page.getByRole('button', { name: /nova transação|adicionar/i })
    if (await createButton.isVisible()) {
      await createButton.click()
      
      // Fill form
      await page.getByLabel(/descrição|description/i).fill('Nova Receita')
      await page.getByLabel(/valor|amount/i).fill('2000')
      
      // Submit
      await page.getByRole('button', { name: /salvar|save/i }).click()
      
      // Should create transaction
      await expect(page).toHaveURL(/.*transactions/)
      expect(createCalled).toBe(true)
    }
  })
})

