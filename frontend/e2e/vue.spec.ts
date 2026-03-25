import { test, expect } from '@playwright/test'

test.describe('JSON Beautifier App', () => {
  test.beforeEach(async ({ page }) => {
    // Mock the beautify endpoint
    await page.route('**/api/v1/beautify', async (route) => {
      const body = route.request().postDataJSON()
      if (body.json === '{bad}') {
        await route.fulfill({
          status: 422,
          contentType: 'application/json',
          body: JSON.stringify({
            code: 'INVALID_JSON',
            message: 'input is not valid JSON',
          }),
        })
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ result: '{\n  "name": "Alice"\n}' }),
        })
      }
    })

    // Mock the minify endpoint
    await page.route('**/api/v1/minify', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ result: '{"name":"Alice"}' }),
      })
    })

    // Mock the validate endpoint
    await page.route('**/api/v1/validate', async (route) => {
      const body = route.request().postDataJSON()
      const valid = body.json !== '{bad}'
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          valid,
          message: valid ? 'JSON is valid' : 'input is not valid JSON',
        }),
      })
    })

    await page.goto('/')
  })

  test('happy path: beautify valid JSON', async ({ page }) => {
    await page.locator('[data-testid="input-area"]').fill('{"name":"Alice"}')
    await page.locator('[data-testid="btn-beautify"]').click()

    const output = await page.locator('[data-testid="output-area"]').inputValue()
    expect(output).toContain('"name"')
    expect(output).toContain('"Alice"')
    expect(page.locator('[data-testid="error-msg"]')).not.toBeVisible()
  })

  test('invalid JSON shows error message', async ({ page }) => {
    await page.locator('[data-testid="input-area"]').fill('{bad}')
    await page.locator('[data-testid="btn-beautify"]').click()

    await expect(page.locator('[data-testid="error-msg"]')).toBeVisible()
    await expect(page.locator('[data-testid="error-msg"]')).toContainText('not valid JSON')
  })
})

