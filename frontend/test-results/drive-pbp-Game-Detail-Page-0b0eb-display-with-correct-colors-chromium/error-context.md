# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: drive-pbp.spec.ts >> Game Detail Page - Drive-Based Play-by-Play >> drive result badges display with correct colors
- Location: tests/drive-pbp.spec.ts:86:3

# Error details

```
Error: Channel closed
```

```
Error: page.waitForSelector: Test ended.
Call log:
  - waiting for locator('h3:has-text("Q1")') to be visible

```

# Page snapshot

```yaml
- generic [ref=e1]:
  - generic [ref=e3]:
    - banner [ref=e4]:
      - generic [ref=e5]:
        - generic [ref=e6]:
          - heading "ATL @ GB" [level=1] [ref=e7]
          - paragraph [ref=e8]: Lambeau Field • 9/25/2026
        - generic [ref=e9]:
          - link "← Back" [ref=e10] [cursor=pointer]:
            - /url: /season/2026/week/3
          - button "Switch to dark mode" [ref=e11] [cursor=pointer]
    - main [ref=e14]:
      - generic [ref=e15]:
        - generic [ref=e17]:
          - generic [ref=e18]:
            - generic [ref=e19]:
              - generic "ATL" [ref=e20]: AT
              - generic [ref=e21]:
                - paragraph [ref=e22]: "35"
                - paragraph [ref=e23]: ATL
            - generic [ref=e24]:
              - generic [ref=e25]: Lambeau Field
              - generic [ref=e26]: 9/25/2026
            - generic [ref=e27]:
              - generic [ref=e28]:
                - paragraph [ref=e29]: "14"
                - paragraph [ref=e30]: GB
              - generic "GB" [ref=e31]
          - paragraph [ref=e33]: Week 3 • 2026 Season
        - generic [ref=e34]:
          - generic [ref=e35]:
            - button "Stats" [ref=e36] [cursor=pointer]
            - button "Play-by-Play" [active] [ref=e37] [cursor=pointer]
            - button "Notes" [ref=e38] [cursor=pointer]
          - generic [ref=e39]: Loading drives...
  - generic [aria-hidden] [ref=e42]: EPA/Play Allowed
```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | test.describe('Game Detail Page - Drive-Based Play-by-Play', () => {
  4   |   const testGameId = '2026_03_ATL_GB';
  5   | 
  6   |   test.beforeEach(async ({ page }) => {
  7   |     await page.goto(`/game/${testGameId}`);
  8   |     await page.waitForLoadState('networkidle');
  9   |   });
  10  | 
  11  |   test('displays game header with team logos and scores', async ({ page }) => {
  12  |     await expect(page.locator('h1:has-text("ATL @ GB")')).toBeVisible({ timeout: 15000 });
  13  |     await expect(page.locator('[aria-label="ATL"]').first()).toBeVisible();
  14  |     await expect(page.locator('[aria-label="GB"]').first()).toBeVisible();
  15  |   });
  16  | 
  17  |   test('displays three tabs: Stats, Play-by-Play, Notes', async ({ page }) => {
  18  |     await expect(page.locator('button:has-text("Stats")')).toBeVisible();
  19  |     await expect(page.locator('button:has-text("Play-by-Play")')).toBeVisible();
  20  |     await expect(page.locator('button:has-text("Notes")')).toBeVisible();
  21  |   });
  22  | 
  23  |   test('shows Stats tab content by default with charts and stat rows', async ({ page }) => {
  24  |     await expect(page.locator('text=Offense Comparison')).toBeVisible({ timeout: 15000 });
  25  |     await expect(page.locator('text=Defense Comparison')).toBeVisible();
  26  |     await expect(page.locator('text=Situational Performance')).toBeVisible();
  27  |   });
  28  | 
  29  |   test('can switch to Play-by-Play tab and see drive-based view', async ({ page }) => {
  30  |     // Click the Play-by-Play tab with force
  31  |     await page.locator('button:has-text("Play-by-Play")').click({ force: true });
  32  |     
  33  |     // Wait for the tab to be active (check the button class)
  34  |     await page.waitForSelector('button:has-text("Play-by-Play").border-nfl-dark', { timeout: 10000 });
  35  |     
  36  |     // Wait for React to re-render and for drive data to load
  37  |     await page.waitForTimeout(1000);
  38  |     
  39  |     // Wait for drive data to load and render
  40  |     await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });
  41  | 
  42  |     // Quarter accordions should be visible with drive counts
  43  |     await expect(page.locator('h3:has-text("Q1")')).toBeVisible({ timeout: 10000 });
  44  |     await expect(page.locator('h3:has-text("Q2")')).toBeVisible();
  45  |     await expect(page.locator('h3:has-text("Q3")')).toBeVisible();
  46  |     await expect(page.locator('h3:has-text("Q4")')).toBeVisible();
  47  |   });
  48  | 
  49  |   test('quarter accordions display team totals', async ({ page }) => {
  50  |     await page.click('button:has-text("Play-by-Play")');
  51  |     // Wait for drive data to load and render
  52  |     await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });
  53  | 
  54  |     // Team totals should show: team name, points, TOP, yards
  55  |     await expect(page.locator('text=ATL').first()).toBeVisible();
  56  |     await expect(page.locator('text=GB').first()).toBeVisible();
  57  |     await expect(page.locator('text=pts').first()).toBeVisible();
  58  |     await expect(page.locator('text=yds').first()).toBeVisible();
  59  |   });
  60  | 
  61  |   test('can expand quarter accordion to see drives', async ({ page }) => {
  62  |     await page.click('button:has-text("Play-by-Play")');
  63  |     // Wait for drive data to load and render
  64  |     await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });
  65  | 
  66  |     // Q1 should be open by default (quarter <= 2)
  67  |     // Drives should be visible inside Q1
  68  |     await expect(page.locator('text=Drive').first()).toBeVisible({ timeout: 10000 });
  69  |   });
  70  | 
  71  |   test('can expand drive item to see plays', async ({ page }) => {
  72  |     await page.click('button:has-text("Play-by-Play")');
  73  |     // Wait for drive data to load and render
  74  |     await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });
  75  | 
  76  |     // Click on first drive to expand
  77  |     const firstDrive = page.locator('button:has-text("Drive")').first();
  78  |     await expect(firstDrive).toBeVisible({ timeout: 10000 });
  79  |     await firstDrive.click();
  80  | 
  81  |     // Play details should be visible
  82  |     await expect(page.locator('text=pass').first()).toBeVisible({ timeout: 10000 });
  83  |     await expect(page.locator('text=EPA:').first()).toBeVisible();
  84  |   });
  85  | 
  86  |   test('drive result badges display with correct colors', async ({ page }) => {
  87  |     await page.click('button:has-text("Play-by-Play")');
  88  |     // Wait for drive data to load and render
> 89  |     await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });
      |                ^ Error: page.waitForSelector: Test ended.
  90  | 
  91  |     // Result badges should be visible (TOUCHDOWN, FIELD_GOAL, PUNT, etc.)
  92  |     await expect(page.locator('text=TOUCHDOWN').first()).toBeVisible({ timeout: 10000 });
  93  |   });
  94  | 
  95  |   test('end play description displays when drive is expanded', async ({ page }) => {
  96  |     await page.click('button:has-text("Play-by-Play")');
  97  |     // Wait for drive data to load and render
  98  |     await page.waitForSelector('h3:has-text("Q1")', { timeout: 30000 });
  99  | 
  100 |     const firstDrive = page.locator('button:has-text("Drive")').first();
  101 |     await firstDrive.click();
  102 | 
  103 |     // END PLAY should be visible
  104 |     await expect(page.locator('text=END PLAY').first()).toBeVisible({ timeout: 10000 });
  105 |   });
  106 | 
  107 |   test('can switch to Notes tab and see game info', async ({ page }) => {
  108 |     await page.click('button:has-text("Notes")');
  109 |     await page.waitForLoadState('networkidle');
  110 | 
  111 |     await expect(page.locator('h3:has-text("Game Information")')).toBeVisible({ timeout: 10000 });
  112 |     await expect(page.locator('h3:has-text("Weather")')).toBeVisible();
  113 |     await expect(page.locator('h3:has-text("Injuries")')).toBeVisible();
  114 |   });
  115 | 
  116 |   test('back button navigates to week page', async ({ page }) => {
  117 |     await page.click('text=← Back');
  118 |     await page.waitForLoadState('networkidle');
  119 | 
  120 |     await expect(page).toHaveURL(/\/season\/2026\/week\/3/);
  121 |   });
  122 | });
  123 | 
  124 | test.describe('Game Detail Page - Drive-Based PBP Error States', () => {
  125 |   test('shows loading state initially', async ({ page }) => {
  126 |     await page.goto('/game/2026_03_ATL_GB');
  127 |     
  128 |     const loadingText = page.locator('text=Loading drives');
  129 |     await expect(loadingText).toBeVisible({ timeout: 1000 });
  130 |   });
  131 | 
  132 |   test('shows no data message for invalid game id', async ({ page }) => {
  133 |     await page.goto('/game/invalid_game_id');
  134 |     await page.waitForLoadState('networkidle');
  135 | 
  136 |     await expect(page.locator('text=No drive data available')).toBeVisible({ timeout: 10000 });
  137 |   });
  138 | });
```