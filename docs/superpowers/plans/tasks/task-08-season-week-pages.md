# Task 8: Season & Week Pages

**Files:**
- Create: `frontend/src/pages/SeasonPage.tsx`
- Create: `frontend/src/pages/WeekPage.tsx`
- Create: `frontend/src/components/game/GameCard.tsx`

**Interfaces:**
- Consumes: `useGames(season, week)`, `useWeeks(season)`
- Produces: Season grid, Week grid with game cards

---

## Steps

- [ ] **Step 1: Implement SeasonPage** — grid of week links

```tsx
// frontend/src/pages/SeasonPage.tsx
export function SeasonPage() {
  const { season } = useParams();
  const seasonNum = parseInt(season || '2024');
  const { data: weeks, isLoading } = useWeeks(seasonNum);

  if (isLoading) return <Loading />;

  return (
    <div>
      <h1 className="text-3xl font-bold text-nfl-dark mb-6">Season {seasonNum}</h1>
      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
        {weeks?.map((week) => (
          <Link key={week} to={`/season/${seasonNum}/week/${week}`} className="card p-6 text-center hover:shadow-md transition-shadow">
            <span className="text-2xl font-bold text-nfl-dark">Week {week}</span>
          </Link>
        ))}
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Implement WeekPage** — list of GameCards

```tsx
// frontend/src/pages/WeekPage.tsx
export function WeekPage() {
  const { season, week } = useParams();
  const { data: games, isLoading } = useGames(parseInt(season!), parseInt(week!));

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-nfl-dark">Week {week}</h1>
          <p className="text-gray-600">Season {season}</p>
        </div>
        <Link to={`/season/${season}`} className="btn-secondary">← Back</Link>
      </div>
      {isLoading ? <Loading /> : games?.map(g => <GameCard key={g.game_id} game={g} />)}
    </div>
  );
}
```

- [ ] **Step 3: Implement GameCard** — expandable row: "Chargers 33 × 10 Broncos [v]" with team logos, scores, click → `/game/{gameId}`

```tsx
// frontend/src/components/game/GameCard.tsx
export function GameCard({ game }: { game: Game }) {
  const [expanded, setExpanded] = useState(false);
  const isHomeWinner = game.home_score > game.away_score;

  return (
    <article className="card">
      <button onClick={() => setExpanded(!expanded)} className="w-full p-4 flex items-center justify-between">
        <div className="flex items-center gap-4 flex-1">
          <TeamLogo team={game.away_team} size="lg" />
          <div className="text-right">
            <p className={`font-bold ${isHomeWinner ? '' : 'text-nfl-light'}`}>{game.away_score}</p>
            <p className="text-sm text-gray-500">{game.away_team}</p>
          </div>
          <span className="text-xl font-bold text-nfl-dark px-4">×</span>
          <div className="text-left">
            <p className={`font-bold ${isHomeWinner ? 'text-nfl-light' : ''}`}>{game.home_score}</p>
            <p className="text-sm text-gray-500">{game.home_team}</p>
          </div>
          <TeamLogo team={game.home_team} size="lg" />
        </div>
        <ChevronIcon className={expanded ? 'rotate-180' : ''} />
      </button>
      {expanded && (
        <div className="border-t bg-gray-50 p-4">
          <Link to={`/game/${game.game_id}`} className="btn-primary w-full">View Game Details</Link>
        </div>
      )}
    </article>
  );
}
```

- [ ] **Step 4: Verify in browser** `npm run dev`

- [ ] **Step 5: Commit**
```bash
git add frontend/src/pages/SeasonPage.tsx frontend/src/pages/WeekPage.tsx frontend/src/components/game/
git commit -m "feat: season and week pages with game cards"
```