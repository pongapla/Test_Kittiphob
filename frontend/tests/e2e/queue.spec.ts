import { test, expect } from '@playwright/test'

test.describe('Queue Management System', () => {

  test('TC-E2E-01: Receive ticket and reset queue', async ({ page }) => {
    await page.goto('/')

    // Reset queue to a known initial state.
    await page.getByTestId('reset-navigation').click()
    await page.getByTestId('reset-button').click()

    await expect(
      page.getByTestId('reset-number')
    ).toHaveText('00')

    // Receive first ticket.
    await page.getByTestId('back-button').click()
    await page.getByTestId('receive-button').click()

    await expect(page).toHaveURL(/\/ticket\/.+/)
    await expect(
      page.getByTestId('queue-number')
    ).toHaveText('A0')

    // Refresh must preserve the same ticket.
    await page.reload()
    await expect(
      page.getByTestId('queue-number')
    ).toHaveText('A0')

    // Receive second ticket.
    await page.getByTestId('back-button').click()
    await page.getByTestId('receive-button').click()

    await expect(
      page.getByTestId('queue-number')
    ).toHaveText('A1')

    // Verify current queue before reset.
    await page.getByTestId('back-button').click()
    await page.getByTestId('reset-navigation').click()

    await expect(
      page.getByTestId('reset-number')
    ).toHaveText('A1')

    // Reset queue.
    await page.getByTestId('reset-button').click()

    await expect(
      page.getByTestId('reset-number')
    ).toHaveText('00')

    // New cycle starts at A0.
    await page.getByTestId('back-button').click()
    await page.getByTestId('receive-button').click()

    await expect(
      page.getByTestId('queue-number')
    ).toHaveText('A0')
  })

})
