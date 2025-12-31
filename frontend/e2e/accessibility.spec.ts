import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

test.describe('Accessibility Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Mock authentication
    await page.goto('/')
    await page.evaluate(() => {
      localStorage.setItem('auth_token', 'mock-token-123')
    })
  })

  test('should have no accessibility violations on login page', async ({ page }) => {
    await page.goto('/login')
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .analyze()
    
    expect(accessibilityScanResults.violations).toEqual([])
  })

  test('should have no accessibility violations on dashboard', async ({ page }) => {
    // Mock API responses
    await page.route('**/api/v1/**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({}),
      })
    })
    
    await page.goto('/')
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .analyze()
    
    expect(accessibilityScanResults.violations).toEqual([])
  })

  test('should have no accessibility violations on transactions page', async ({ page }) => {
    // Mock API responses
    await page.route('**/api/v1/transactions**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          transactions: [],
          count: 0,
        }),
      })
    })
    
    await page.goto('/transactions')
    
    const accessibilityScanResults = await new AxeBuilder({ page })
      .analyze()
    
    expect(accessibilityScanResults.violations).toEqual([])
  })

  test('should have proper heading hierarchy', async ({ page }) => {
    await page.goto('/login')
    
    // Check for h1
    const h1 = await page.locator('h1').count()
    expect(h1).toBeGreaterThanOrEqual(1)
  })

  test('should have proper form labels', async ({ page }) => {
    await page.goto('/login')
    
    // Check for form labels
    const emailLabel = page.getByLabel(/email/i)
    const passwordLabel = page.getByLabel(/password|senha/i)
    
    await expect(emailLabel).toBeVisible()
    await expect(passwordLabel).toBeVisible()
  })

  test('should have proper button labels', async ({ page }) => {
    await page.goto('/login')
    
    // Check for button with accessible name
    const submitButton = page.getByRole('button', { name: /entrar|login/i })
    await expect(submitButton).toBeVisible()
  })
})

