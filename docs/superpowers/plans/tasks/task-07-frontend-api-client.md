# Task 7: Frontend Scaffolding & API Client

**Files:**
- Create: `frontend/src/api/client.ts`
- Create: `frontend/src/api/types.ts`
- Create: `frontend/src/hooks/useGames.ts`
- Create: `frontend/src/hooks/useGameStats.ts`
- Create: `frontend/src/hooks/usePlays.ts`

**Interfaces:**
- Consumes: Vite proxy config (`/api` → `http://localhost:8080`)
- Produces: Typed API functions, React Query hooks

---

## Steps

- [ ] **Step 1: Define TypeScript types** matching Go models (`frontend/src/api/types.ts`)

```typescript
export interface Game {
  game_id: string;
  season: number;
  week: number;
  home_team: string;
  away_team: string;
  home_score: number;
  away_score: number;
  game_date: string;
  venue: string;
}

export interface TeamComparisonStats {
  home_team: string;
  away_team: string;
  home_offense: TeamUnitStats;
  home_defense: TeamUnitStats;
  away_offense: TeamUnitStats;
  away_defense: TeamUnitStats;
}

export interface TeamUnitStats {
  plays: number;
  epa_per_play: number;
  success_rate: number;
  pass_rate: number;
  rush_epa: number;
  pass_epa: number;
  red_zone_pct: number;
  third_down_pct: number;
  fourth_down_pct: number;
  time_of_possession: string;
  explosive_plays: number;
  turnovers: number;
}

export interface PlayDetail {
  play_id: number;
  quarter: number;
  down: number;
  distance: number;
  yardline: number;
  play_type: string;
  description: string;
  epa: number;
  passer: string;
  receiver: string;
  rusher: string;
  tackler: string;
  is_scoring_play: boolean;
  drive_id: number;
  drive_play_count: number;
  drive_yards: number;
  drive_result: string;
  drive_top: string;
}

export interface APIResponse<T> {
  data: T;
  meta: { cached: boolean; season: number };
}
```

- [ ] **Step 2: Implement ky-based API client** (`frontend/src/api/client.ts`)

```typescript
import ky from 'ky';

const api = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 30000,
  retry: 1,
});

export const getSeasons = () => api.get('api/seasons').json<number[]>();
export const getWeeks = (season: number) => api.get(`api/seasons/${season}/weeks`).json<number[]>();
export const getGames = (season: number, week: number) => api.get(`api/seasons/${season}/weeks/${week}/games`).json<Game[]>();
export const getGame = (gameId: string) => api.get(`api/games/${gameId}`).json<Game>();
export const getGameStats = (gameId: string) => api.get(`api/games/${gameId}/stats`).json<APIResponse<TeamComparisonStats>>();
export const getPlays = (gameId: string, params?: { quarter?: number; drive?: number }) => 
  api.get(`api/games/${gameId}/plays`, { searchParams: params }).json<APIResponse<PlayDetail[]>>();
export const getGameNotes = (gameId: string) => api.get(`api/games/${gameId}/notes`).json<APIResponse<GameNotes>>();
```

- [ ] **Step 3: Create React Query hooks** (`frontend/src/hooks/useGames.ts`, etc.)

```typescript
// frontend/src/hooks/useGames.ts
export function useGames(season: number, week: number) {
  return useQuery({
    queryKey: ['games', season, week],
    queryFn: () => api.getGames(season, week),
    enabled: !!season && !!week,
  });
}

export function useWeeks(season: number) {
  return useQuery({
    queryKey: ['weeks', season],
    queryFn: () => api.getWeeks(season),
    enabled: !!season,
  });
}
```

- [ ] **Step 4: Configure QueryClient** in `main.tsx` (staleTime 5min, retry 1)

- [ ] **Step 5: Verify types compile** `npm run build`

- [ ] **Step 6: Commit**
```bash
git add frontend/src/api/ frontend/src/hooks/
git commit -m "feat: frontend api client and react query hooks"
```