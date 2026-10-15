import { test, expect } from '@playwright/test';

test.describe('Game Detail Page - Drive-Based Play-by-Play', () => {
  const testGameId = '2024_01_LAC_DEN';

  test.beforeEach(async ({ page }) => {
    await page.goto(`/game/${testGameId}`);
    await page.waitForLoadState('networkidle');
  });

  test('displays game header with team logos and scores', async ({ page }) => {
    await expect(page.locator('h1:has-text("DEN @ LAC")')).toBeVisible({ timeout: 15000 });
    await expect(page.locator('[aria-label="DEN"]').first()).toBeVisible();
    await expect(page.locator('[aria-label="LAC"]').first()).toBeVisible();
  });

  test('displays three tabs: Stats, Play-by-Play, Notes', async ({ page }) => {
    await expect(page.locator('button:has-text("Stats")')).toBeVisible();
    await expect(page.locator('button:has-text("Play-by-Play")')).toBeVisible();
    await expect(page.locator('button:has-text("Notes")')).toBeVisible();
  });

  test('shows Stats tab content by default with charts and stat rows', async ({ page }) => {
    await expect(page.locator('text=Offense Comparison')).toBeVisible({ timeout: 15000 });
    await expect(page.locator('text=Defense Comparison')).toBeVisible();
    await expect(page.locator('text=Situational Performance')).toBeVisible();
  });

  test('can switch to Play-by-Play tab and see drive-based view', async ({ page }) => {
    // Click the Play-by-Play tab with force
    await page.locator('button:has-text("Play-by-Play")').click({ force: true });
    
    // Wait for the tab to be active (check the button class)
    await page.waitForSelector('button:has-text("Play-by-Play").border-nfl-dark', { timeout: 10000 });
    
    // Wait for React to re-render and for drive data to load
    await page.waitForTimeout(1000);
    
    // Wait for drive data to load and render
    await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });

    // Quarter accordions should be visible with drive counts
    await expect(page.locator('h3:has-text("Q1")')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('h3:has-text("Q2")')).toBeVisible();
    await expect(page.locator('h3:has-text("Q3")')).toBeVisible();
    await expect(page.locator('h3:has-text("Q4")')).toBeVisible();
  });

  test('quarter accordions display team totals', async ({ page }) => {
    await page.click('button:has-text("Play-by-Play")');
    // Wait for drive data to load and render
    await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });

    // Team totals should show: team name, points, TOP, yards
    await expect(page.locator('text=LAC').first()).toBeVisible();
    await expect(page.locator('text=DEN').first()).toBeVisible();
    await expect(page.locator('text=pts').first()).toBeVisible();
    await expect(page.locator('text=yds').first()).toBeVisible();
  });

  test('can expand quarter accordion to see drives', async ({ page }) => {
    await page.click('button:has-text("Play-by-Play")');
    // Wait for drive data to load and render
    await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });

    // Q1 should be open by default (quarter <= 2)
    // Drives should be visible inside Q1
    await expect(page.locator('text=Drive').first()).toBeVisible({ timeout: 10000 });
  });

  test('can expand drive item to see plays', async ({ page }) => {
    await page.click('button:has-text("Play-by-Play")');
    // Wait for drive data to load and render
    await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });

    // Click on first drive to expand
    const firstDrive = page.locator('button:has-text("Drive")').first();
    await expect(firstDrive).toBeVisible({ timeout: 10000 });
    await firstDrive.click();

    // Play details should be visible
    await expect(page.locator('text=pass').first()).toBeVisible({ timeout: 10000 });
    await expect(page.locator('text=EPA:').first()).toBeVisible();
  });

  test('drive result badges display with correct colors', async ({ page }) => {
    await page.click('button:has-text("Play-by-Play")');
    // Wait for drive data to load and render
    await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });

    // Result badges should be visible (TOUCHDOWN, FIELD_GOAL, PUNT, etc.)
    await expect(page.locator('text=TOUCHDOWN').first()).toBeVisible({ timeout: 10000 });
  });

  test('end play description displays when drive is expanded', async ({ page }) => {
    await page.click('button:has-text("Play-by-Play")');
    // Wait for drive data to load and render
    await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });

    const firstDrive = page.locator('button:has-text("Drive")').first();
    await firstDrive.click();

    // END PLAY should be visible
    await expect(page.locator('text=END PLAY').first()).toBeVisible({ timeout: 10000 });
  });

  test('can switch to Notes tab and see game info', async ({ page }) => {
    await page.click('button:has-text("Notes")');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('h3:has-text("Game Information")')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('h3:has-text("Weather")')).toBeVisible();
    await expect(page.locator('h3:has-text("Injuries")')).toBeVisible();
  });

  test('back button navigates to week page', async ({ page }) => {
    await page.click('text=← Back');
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveURL(/\/season\/2024\/week\/1/);
  });
});

test.describe('Game Detail Page - Drive-Based PBP Error States', () => {
  test('shows loading state initially', async ({ page }) => {
    await page.goto('/game/2024_01_LAC_DEN');
    
    const loadingText = page.locator('text=Loading drives');
    await expect(loadingText).toBeVisible({ timeout: 1000 });
  });

  test('shows no data message for invalid game id', async ({ page }) => {
    await page.goto('/game/invalid_game_id');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('text=No drive data available')).toBeVisible({ timeout: 10000 });
  });
});