import { test, expect } from '@playwright/test';

test.describe('Game Detail Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/season/2024');
    await page.waitForLoadState('networkidle');
  });

  test('navigates to game detail page when clicking on a game card', async ({ page }) => {
    const weekLink = page.locator('a[href*="/week/1"]').first();
    await expect(weekLink).toBeVisible({ timeout: 10000 });
    await weekLink.click();
    await page.waitForLoadState('networkidle');

    const gameCard = page.locator('[data-testid="game-card"]').first();
    await expect(gameCard).toBeVisible({ timeout: 10000 });

    const gameId = await gameCard.getAttribute('data-game-id');
    // Expand the card first, then click "View Game Details"
    await gameCard.locator('button[aria-expanded="false"]').click();
    await expect(gameCard.locator('button:has-text("View Game Details")')).toBeVisible({ timeout: 5000 });
    await gameCard.locator('button:has-text("View Game Details")').click();
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveURL(new RegExp(`/game/${gameId}`));
  });

  test('displays game header with team logos and scores', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    // Use more specific locators to avoid strict mode violations
    await expect(page.locator('h1:has-text("LV @ LAC")')).toBeVisible({ timeout: 15000 });
    await expect(page.locator('[aria-label="LV"]').first()).toBeVisible();
    await expect(page.locator('[aria-label="LAC"]').first()).toBeVisible();
  });

  test('displays three tabs: Stats, Play-by-Play, Notes', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('button:has-text("Stats")')).toBeVisible();
    await expect(page.locator('button:has-text("Play-by-Play")')).toBeVisible();
    await expect(page.locator('button:has-text("Notes")')).toBeVisible();
  });

  test('shows Stats tab content by default with charts and stat rows', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('text=Offense Comparison')).toBeVisible({ timeout: 15000 });
    await expect(page.locator('text=Defense Comparison')).toBeVisible();
    await expect(page.locator('text=Situational Performance')).toBeVisible();
  });

  test('displays TeamComparisonCharts with bar charts', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    // Use heading locators to avoid strict mode violations (multiple stat rows have same text)
    await expect(page.locator('h3:has-text("Offense Comparison")')).toBeVisible({ timeout: 15000 });
    await expect(page.locator('h4:has-text("Key Offensive Metrics")')).toBeVisible();
    await expect(page.locator('h3:has-text("Defense Comparison")')).toBeVisible();
    await expect(page.locator('h4:has-text("Key Defensive Metrics")')).toBeVisible();
  });

  test('displays radar charts for situational performance', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    // Use heading locator for situational performance section
    await expect(page.locator('h3:has-text("Situational Performance")')).toBeVisible({ timeout: 15000 });
    // The radar charts have team names as legends
    await expect(page.locator('text=LV Offense')).toBeVisible();
    await expect(page.locator('text=LAC Offense')).toBeVisible();
  });

  test('shows stat rows for each team unit', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('text=LV Offense')).toBeVisible({ timeout: 15000 });
    await expect(page.locator('text=LAC Offense')).toBeVisible();
    await expect(page.locator('text=LV Defense')).toBeVisible();
    await expect(page.locator('text=LAC Defense')).toBeVisible();
  });

  test('can switch to Play-by-Play tab and see plays grouped by quarter', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    // Plays are loaded on initial page load, just wait for tab switch
    await page.click('button:has-text("Play-by-Play")');
    // Wait longer for component to mount and render
    await page.waitForSelector('h4:has-text("1Q")', { timeout: 30000 });

    // Quarter headings include play count: "1Q (N plays)"
    await expect(page.locator('h4:has-text("1Q")')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('h4:has-text("2Q")')).toBeVisible();
    await expect(page.locator('h4:has-text("3Q")')).toBeVisible();
    await expect(page.locator('h4:has-text("4Q")')).toBeVisible();
  });

  test('displays play details in Play-by-Play tab', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await page.click('button:has-text("Play-by-Play")');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('text=pass').first()).toBeVisible({ timeout: 10000 });
    await expect(page.locator('text=rush').first()).toBeVisible();
  });

  test('can switch to Notes tab and see game info, weather, injuries', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await page.click('button:has-text("Notes")');
    await page.waitForLoadState('networkidle');

    // Use heading locators
    await expect(page.locator('h3:has-text("Game Information")')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('h3:has-text("Weather")')).toBeVisible();
    await expect(page.locator('h3:has-text("Injuries")')).toBeVisible();
  });

  test('displays weather information in Notes tab', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await page.click('button:has-text("Notes")');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('text=Temperature')).toBeVisible({ timeout: 10000 });
    await expect(page.locator('text=Conditions')).toBeVisible();
    await expect(page.locator('text=Wind')).toBeVisible();
    await expect(page.locator('text=Humidity')).toBeVisible();
  });

  test('back button navigates to week page', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    await page.waitForLoadState('networkidle');

    await page.click('text=← Back');
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveURL(/\/season\/2024\/week\/1/);
  });
});

test.describe('Game Detail Page - Error States', () => {
  test('shows loading state initially', async ({ page }) => {
    await page.goto('/game/2024_01_LV_LAC');
    
    const loadingText = page.locator('text=Loading game data');
    await expect(loadingText).toBeVisible({ timeout: 1000 });
  });

  test('shows game not found for invalid game id', async ({ page }) => {
    await page.goto('/game/invalid_game_id');
    await page.waitForLoadState('networkidle');

    await expect(page.locator('text=Game not found')).toBeVisible({ timeout: 10000 });
  });
});