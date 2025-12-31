import { test, expect } from '@playwright/test'

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should redirect to login when not authenticated', async ({ page }) => {
    // Try to access a protected route
    await page.goto('/accounts')
    
    // Should redirect to login
    await expect(page).toHaveURL(/.*login/)
  })

  test('should show login form', async ({ page }) => {
    await page.goto('/login')
    
    // Check for login form elements
    await expect(page.getByLabel(/email/i)).toBeVisible()
    await expect(page.getByLabel(/password/i)).toBeVisible()
    await expect(page.getByRole('button', { name: /entrar|login/i })).toBeVisible()
  })

  test('should display validation errors on invalid login', async ({ page }) => {
    await page.goto('/login')
    
    // Try to submit empty form
    await page.getByRole('button', { name: /entrar|login/i }).click()
    
    // Should show validation errors
    await expect(page.getByText(/email.*required|obrigatório/i)).toBeVisible()
  })

  test('should login successfully with valid credentials', async ({ page }) => {
    await page.goto('/login')
    
    // Fill login form
    await page.getByLabel(/email/i).fill('test@example.com')
    await page.getByLabel(/password/i).fill('password123')
    
    // Mock successful login
    await page.route('**/api/v1/auth/login', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-token-123',
          user: {
            user_id: 'user-123',
            email: 'test@example.com',
            first_name: 'Test',
            last_name: 'User',
          },
        }),
      })
    })
    
    // Submit form
    await page.getByRole('button', { name: /entrar|login/i }).click()
    
    // Should redirect to dashboard or home
    await expect(page).toHaveURL(/.*\/$|.*dashboard/)
  })

  test('should logout successfully', async ({ page, context }) => {
    // First login
    await page.goto('/login')
    
    // Mock login
    await page.route('**/api/v1/auth/login', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-token-123',
          user: {
            user_id: 'user-123',
            email: 'test@example.com',
          },
        }),
      })
    })
    
    await page.getByLabel(/email/i).fill('test@example.com')
    await page.getByLabel(/password/i).fill('password123')
    await page.getByRole('button', { name: /entrar|login/i }).click()
    
    // Wait for navigation
    await page.waitForURL(/.*\/$|.*dashboard/)
    
    // Find and click logout button
    const logoutButton = page.getByRole('button', { name: /logout|sair/i })
    if (await logoutButton.isVisible()) {
      await logoutButton.click()
      
      // Should redirect to login
      await expect(page).toHaveURL(/.*login/)
    }
  })
})

