# Task 14: Polish & Edge Cases

**Files:**
- Modify: various

---

## Steps

- [ ] **Step 1: Error boundaries** in React for failed queries
  - Wrap page components in `ErrorBoundary`
  - Show user-friendly error with retry button
  - Log errors to console/sentry

- [ ] **Step 2: Loading skeletons** for all async components
  - `GameCardSkeleton`, `StatsSkeleton`, `PlayByPlaySkeleton`, `NotesSkeleton`
  - Use `@tanstack/react-query` `isLoading` state

- [ ] **Step 3: Empty states** (no games, no plays, no data)
  - Season: "No weeks available"
  - Week: "No games this week"
  - Game Stats: "Stats unavailable"
  - PBP: "No plays recorded"
  - Notes: "No notes available"

- [ ] **Step 4: Responsive design**
  - Mobile week grid (1 col), game card (stacked)
  - Tablet: 2-3 col grids
  - Desktop: full layout
  - Test at 375px, 768px, 1024px, 1440px

- [ ] **Step 5: Accessibility**
  - ARIA labels on all interactive elements
  - Keyboard navigation (tab order, focus visible)
  - Color contrast (WCAG AA)
  - Semantic HTML (header, main, nav, article, section)
  - Screen reader announcements for live regions

- [ ] **Step 6: Performance**
  - `React.memo` for `GameCard`, `PlayRow`, `StatRow`
  - `useMemo` for chart data transformations
  - Virtualized list for PBP (react-window) if >100 plays
  - Code splitting: lazy load GamePage components
  - Bundle analysis: `npm run build && npx vite-bundle-analyzer`

- [ ] **Step 7: Run full test suite**
```bash
cd backend && go test ./... && go vet ./... && golangci-lint run
cd ../frontend && npm run test && npm run lint && npm run build
```

- [ ] **Step 8: Backend polish**
  - Add request ID middleware for tracing
  - Structured JSON logging with slog
  - Health check endpoint with dependency checks (Redis, GitHub)
  - Graceful shutdown on SIGTERM
  - Rate limit handling with exponential backoff for GitHub API

- [ ] **Step 9: Create feature branch**
```bash
git checkout -b feat/polish-edge-cases
```

- [ ] **Step 10: Commit**
```bash
git add -A
git commit -m "feat: polish - error boundaries, skeletons, a11y, performance, empty states"
```