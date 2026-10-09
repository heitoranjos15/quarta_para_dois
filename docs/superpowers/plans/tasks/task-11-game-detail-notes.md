# Task 11: Game Detail Page (Notes Tab)

**Files:**
- Create: `frontend/src/components/notes/GameNotes.tsx`
- Create: `frontend/src/components/notes/Weather.tsx`
- Create: `frontend/src/components/notes/Injuries.tsx`
- Create: `frontend/src/components/notes/VegasLines.tsx`

**Interfaces:**
- Consumes: `useGameNotes(gameId)`
- Produces: Game metadata, weather, injuries, vegas lines

---

## Steps

- [ ] **Step 1: Define GameNotes type** (`frontend/src/api/types.ts`)

```typescript
export interface GameNotes {
  game_info: {
    stadium: string;
    surface: string;
    roof: string;
    attendance: number;
    duration: string;
    officials: string[];
  };
  weather: {
    temperature: number;
    wind: number;
    humidity: number;
    conditions: string;
  } | null;
  injuries: {
    team: string;
    player: string;
    position: string;
    status: string;
    detail: string;
  }[];
  vegas: {
    spread: number;
    total: number;
    home_moneyline: number;
    away_moneyline: number;
  } | null;
}
```

- [ ] **Step 2: Implement Weather component** (`frontend/src/components/notes/Weather.tsx`)

```tsx
export function Weather({ weather }: { weather: GameNotes['weather'] }) {
  if (!weather) return <div className="text-gray-500">Weather data unavailable</div>;
  return (
    <div className="grid grid-cols-4 gap-4 text-center">
      <div className="p-4 bg-blue-50 rounded-lg">
        <p className="text-3xl font-bold text-blue-600">{weather.temperature}°F</p>
        <p className="text-sm text-gray-500">Temperature</p>
      </div>
      <div className="p-4 bg-green-50 rounded-lg">
        <p className="text-3xl font-bold text-green-600">{weather.wind} mph</p>
        <p className="text-sm text-gray-500">Wind</p>
      </div>
      <div className="p-4 bg-yellow-50 rounded-lg">
        <p className="text-3xl font-bold text-yellow-600">{weather.humidity}%</p>
        <p className="text-sm text-gray-500">Humidity</p>
      </div>
      <div className="p-4 bg-gray-50 rounded-lg">
        <p className="text-sm font-medium text-gray-700">{weather.conditions}</p>
        <p className="text-xs text-gray-500">Conditions</p>
      </div>
    </div>
  );
}
```

- [ ] **Step 3: Implement Injuries component** (`frontend/src/components/notes/Injuries.tsx`)

```tsx
export function Injuries({ injuries }: { injuries: GameNotes['injuries'] }) {
  if (!injuries?.length) return <div className="text-gray-500">No injury reports</div>;
  return (
    <div className="space-y-2">
      {injuries.map((injury, i) => (
        <div key={i} className="p-3 bg-white border rounded-lg">
          <div className="flex items-center justify-between">
            <span className="font-medium">{injury.player} ({injury.position})</span>
            <span className="text-sm text-gray-500">{injury.team}</span>
          </div>
          <div className="flex items-center gap-2 mt-1">
            <span className={`px-2 py-0.5 rounded text-xs ${getStatusColor(injury.status)}`}>{injury.status}</span>
            <span className="text-sm text-gray-600">{injury.detail}</span>
          </div>
        </div>
      ))}
    </div>
  );
}
```

- [ ] **Step 4: Implement VegasLines component** (`frontend/src/components/notes/VegasLines.tsx`)

```tsx
export function VegasLines({ vegas }: { vegas: GameNotes['vegas'] }) {
  if (!vegas) return <div className="text-gray-500">Vegas lines unavailable</div>;
  return (
    <div className="grid grid-cols-2 gap-4">
      <StatCard label="Spread" value={vegas.spread > 0 ? `+${vegas.spread}` : vegas.spread} />
      <StatCard label="Total" value={vegas.total} />
      <StatCard label="Home ML" value={vegas.home_moneyline > 0 ? `+${vegas.home_moneyline}` : vegas.home_moneyline} />
      <StatCard label="Away ML" value={vegas.away_moneyline > 0 ? `+${vegas.away_moneyline}` : vegas.away_moneyline} />
    </div>
  );
}
```

- [ ] **Step 5: Implement GameNotes panel** (`frontend/src/components/notes/GameNotes.tsx`) combining all

- [ ] **Step 6: Wire into GamePage Notes tab**

- [ ] **Step 7: Create feature branch**
```bash
git checkout -b feat/game-detail-notes
```

- [ ] **Step 8: Commit**
```bash
git add frontend/src/components/notes/
git commit -m "feat: game notes tab with weather, injuries, vegas"
```