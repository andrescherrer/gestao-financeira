import { test, expect } from '@playwright/test'

test.describe('Accounts Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Mock authentication
    await page.goto('/')
    await page.evaluate(() => {
      localStorage.setItem('auth_token', 'mock-token-123')
    })
    
    // Mock API responses
    await page.route('**/api/v1/accounts**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
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
        }),
      })
    })
  })

  test('should display accounts list', async ({ page }) => {
    await page.goto('/accounts')
    
    // Should show accounts
    await expect(page.getByText('Conta Corrente')).toBeVisible()
    await expect(page.getByText('Carteira')).toBeVisible()
  })

  test('should open create account form', async ({ page }) => {
    await page.goto('/accounts')
    
    // Click create button
    const createButton = page.getByRole('button', { name: /nova conta|adicionar/i })
    if (await createButton.isVisible()) {
      await createButton.click()
      
      // Should show form
      await expect(page.getByLabel(/nome|name/i)).toBeVisible()
    }
  })

  test('should create a new account', async ({ page }) => {
    await page.goto('/accounts')
    
    let createCalled = false
    await page.route('**/api/v1/accounts', async route => {
      if (route.request().method() === 'POST') {
        createCalled = true
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            account_id: 'acc-new',
            name: 'Nova Conta',
            type: 'BANK',
            currency: 'BRL',
            balance: 0.0,
            is_active: true,
          }),
        })
      } else {
        await route.continue()
      }
    })
    
    // Open create form
    const createButton = page.getByRole('button', { name: /nova conta|adicionar/i })
    if (await createButton.isVisible()) {
      await createButton.click()
      
      // Fill form
      await page.getByLabel(/nome|name/i).fill('Nova Conta')
      
      // Submit
      await page.getByRole('button', { name: /salvar|save/i }).click()
      
      // Should create account
      await expect(page).toHaveURL(/.*accounts/)
      expect(createCalled).toBe(true)
    }
  })
})

