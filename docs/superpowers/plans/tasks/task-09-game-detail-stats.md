# Task 9: Game Detail Page (Stats Tab)

**Files:**
- Create: `frontend/src/pages/GamePage.tsx`
- Create: `frontend/src/components/stats/TeamComparisonCharts.tsx`
- Create: `frontend/src/components/stats/StatRow.tsx`
- Create: `frontend/src/components/ui/Tabs.tsx`

**Interfaces:**
- Consumes: `useGameStats(gameId)`, `useGameInfo(gameId)`
- Produces: Game header + tabs (Stats | Play-by-Play | Notes)

---

## Steps

- [ ] **Step 1: Implement Tabs component** (`frontend/src/components/ui/Tabs.tsx`)

```tsx
export function Tabs({ tabs, defaultTab }: { tabs: {id: string, label: string}[], defaultTab: string }) {
  const [active, setActive] = useState(defaultTab);
  return (
    <div>
      <div className="flex border-b border-gray-200">
        {tabs.map(t => (
          <button key={t.id} onClick={() => setActive(t.id)} className={`px-4 py-2 border-b-2 font-medium transition-colors ${active === t.id ? 'border-nfl-dark text-nfl-dark' : 'border-transparent text-gray-500 hover:text-gray-700'}`}>
            {t.label}
          </button>
        ))}
      </div>
      <div className="py-4">{/* render active tab panel */}</div>
    </div>
  );
}
```

- [ ] **Step 2: Implement GamePage layout** with tabs (`frontend/src/pages/GamePage.tsx`)

```tsx
export function GamePage() {
  const { gameId } = useParams();
  const { data: game } = useGameInfo(gameId!);
  const { data: stats } = useGameStats(gameId!);
  const { data: plays } = usePlays(gameId!);
  const { data: notes } = useGameNotes(gameId!);

  return (
    <div>
      <GameHeader game={game} />
      <Tabs tabs={[
        {id: 'stats', label: 'Stats'},
        {id: 'pbp', label: 'Play-by-Play'},
        {id: 'notes', label: 'Notes'},
      ]} defaultTab="stats">
        <TabPanel id="stats"><GameStats stats={stats?.data} /></TabPanel>
        <TabPanel id="pbp"><PlayByPlay plays={plays?.data} /></TabPanel>
        <TabPanel id="notes"><GameNotesPanel notes={notes?.data} /></TabPanel>
      </Tabs>
    </div>
  );
}
```

- [ ] **Step 3: Implement TeamComparisonCharts** (`frontend/src/components/stats/TeamComparisonCharts.tsx`) using Recharts

Charts needed:
- Bar chart: EPA/play (Home Off vs Away Off, Home Def vs Away Def)
- Bar chart: Success Rate
- Grouped bar: Pass vs Rush EPA
- Radial/radar chart: Red Zone, 3rd Down, 4th Down, Explosive Plays

```tsx
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, RadarChart, Radar, PolarGrid, PolarAngleAxis, PolarRadiusAxis } from 'recharts';

export function TeamComparisonCharts({ stats }: { stats: TeamComparisonStats }) {
  const offenseData = [
    { metric: 'EPA/Play', home: stats.home_offense.epa_per_play, away: stats.away_offense.epa_per_play },
    { metric: 'Success Rate', home: stats.home_offense.success_rate * 100, away: stats.away_offense.success_rate * 100 },
    { metric: 'Pass EPA', home: stats.home_offense.pass_epa, away: stats.away_offense.pass_epa },
    { metric: 'Rush EPA', home: stats.home_offense.rush_epa, away: stats.away_offense.rush_epa },
  ];
  // ... render BarChart and RadarChart
}
```

- [ ] **Step 4: Implement StatRow** component for key-value display

```tsx
export function StatRow({ label, home, away, unit = '' }: { label: string; home: number; away: number; unit?: string }) {
  return (
    <div className="flex justify-between items-center py-2 border-b border-gray-100">
      <span className="text-gray-600">{label}</span>
      <div className="flex items-center gap-4">
        <span className="font-medium text-nfl-dark">{home}{unit}</span>
        <span className="text-gray-400">vs</span>
        <span className="font-medium text-nfl-light">{away}{unit}</span>
      </div>
    </div>
  );
}
```

- [ ] **Step 5: Verify with real API** (start backend, visit `/game/2024_01_LAC_DEN`)

- [ ] **Step 6: Commit**
```bash
git add frontend/src/pages/GamePage.tsx frontend/src/components/stats/ frontend/src/components/ui/Tabs.tsx
git commit -m "feat: game detail page with stats tab and charts"
```