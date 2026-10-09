# Task 10: Game Detail Page (Play-by-Play Tab)

**Files:**
- Create: `frontend/src/components/pbp/PlayByPlayTable.tsx`
- Create: `frontend/src/components/pbp/QuarterAccordion.tsx`
- Create: `frontend/src/components/pbp/PlayRow.tsx`

**Interfaces:**
- Consumes: `usePlays(gameId, {quarter?, drive?})`
- Produces: Quarter accordions with play rows

---

## Play Row Format

> "1 & 10 pass Herbert (QB) to Wilson (TE) 15 yards (EPA: +1.2) — tackled by Hurts (DB)"

---

## Steps

- [ ] **Step 1: Implement QuarterAccordion** (`frontend/src/components/pbp/QuarterAccordion.tsx`)

```tsx
export function QuarterAccordion({ quarter, plays }: { quarter: number; plays: PlayDetail[] }) {
  const [open, setOpen] = useState(quarter <= 2);

  return (
    <details className="group" open={open}>
      <summary className="flex items-center justify-between p-3 bg-gray-50 border-b cursor-pointer list-none">
        <h3 className="font-semibold text-nfl-dark">Q{quarter} <span className="text-sm font-normal text-gray-500">({plays.length} plays)</span></h3>
        <ChevronIcon className="transition-transform group-open:rotate-180" />
      </summary>
      <div className="p-2 space-y-1">
        {plays.map(play => <PlayRow key={play.play_id} play={play} />)}
      </div>
    </details>
  );
}
```

- [ ] **Step 2: Implement PlayRow** (`frontend/src/components/pbp/PlayRow.tsx`)

```tsx
export function PlayRow({ play }: { play: PlayDetail }) {
  const epaColor = play.epa > 0 ? 'text-green-600' : play.epa < 0 ? 'text-red-600' : 'text-gray-500';
  const playTypeIcon = getPlayTypeIcon(play.play_type);

  return (
    <div className="p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors">
      <div className="flex items-start gap-3">
        <span className="w-16 text-xs text-gray-500 font-mono">{play.down}&{play.distance}</span>
        <span className="w-8">{playTypeIcon}</span>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium">{play.description}</p>
          <p className="text-xs text-gray-500">
            {play.passer && `QB: ${play.passer}`}
            {play.receiver && ` → ${play.receiver}`}
            {play.rusher && ` RB: ${play.rusher}`}
            {play.tackler && ` — tackled by ${play.tackler}`}
          </p>
        </div>
        <span className={`${epaColor} font-mono font-medium whitespace-nowrap`}>EPA: {play.epa > 0 ? '+' : ''}{play.epa.toFixed(1)}</span>
      </div>
      {play.is_scoring_play && <span className="text-xs bg-nfl-gold text-nfl-dark px-2 py-0.5 rounded">SCORING PLAY</span>}
    </div>
  );
}
```

- [ ] **Step 3: Implement PlayByPlayTable** with quarter filter tabs (`frontend/src/components/pbp/PlayByPlayTable.tsx`)

```tsx
export function PlayByPlayTable({ plays }: { plays: PlayDetail[] }) {
  const [quarterFilter, setQuarterFilter] = useState<number | 'all'>('all');
  const quarters = [1,2,3,4].filter(q => plays.some(p => p.quarter === q));
  const filtered = quarterFilter === 'all' ? plays : plays.filter(p => p.quarter === quarterFilter);

  return (
    <div>
      <div className="flex gap-2 mb-4 border-b pb-2">
        <button onClick={() => setQuarterFilter('all')} className={`px-3 py-1 rounded ${quarterFilter === 'all' ? 'bg-nfl-dark text-white' : 'bg-gray-100'}`}>All</button>
        {quarters.map(q => (
          <button key={q} onClick={() => setQuarterFilter(q)} className={`px-3 py-1 rounded ${quarterFilter === q ? 'bg-nfl-medium text-white' : 'bg-gray-100'}`}>Q{q}</button>
        ))}
      </div>
      <div className="space-y-2">
        {quarterFilter === 'all' 
          ? quarters.map(q => <QuarterAccordion key={q} quarter={q} plays={plays.filter(p => p.quarter === q)} />)
          : <QuarterAccordion quarter={quarterFilter} plays={filtered} />
        }
      </div>
    </div>
  );
}
```

- [ ] **Step 4: Verify with real data**

- [ ] **Step 5: Commit**
```bash
git add frontend/src/components/pbp/
git commit -m "feat: play-by-play tab with quarter accordions"
```