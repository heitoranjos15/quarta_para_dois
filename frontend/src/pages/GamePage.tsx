import { useParams, Link } from 'react-router-dom';
import { useGame, useGameNotes } from '../hooks/useGames';
import { useGameStats } from '../hooks/useGameStats';
import { useDrives } from '../hooks/useDrives';
import { ThemeToggle } from '../components/ui/ThemeToggle';
import { TeamLogo } from '../components/ui/TeamLogo';
import { Tabs, TabPanel } from '../components/ui/Tabs';
import { TeamComparisonCharts } from '../components/stats/TeamComparisonCharts';
import { StatRow } from '../components/stats/StatRow';
import { DriveList } from '../components/pbp/DriveList';
import type { Game, TeamComparisonStats, TeamUnitStats, GameNotes } from '../api/types';

const formatKickoff = (kickoff: string) => (kickoff ? new Date(kickoff).toLocaleString() : '');
const mmss = (seconds: number) =>
  `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, '0')}`;
const fixed2 = (n: number) => n.toFixed(2);
const percent = (n: number) => `${Math.round(n)}%`;

export function GamePage() {
  const params = useParams();
  const gameId = params.gameId || '';

  const { data: game, isLoading: gameLoading } = useGame(gameId);
  const { data: stats, isLoading: statsLoading } = useGameStats(gameId);
  const { data: notes } = useGameNotes(gameId);

  if (gameLoading || statsLoading) {
    return <div className="p-8 text-center text-[var(--text-secondary)]">Loading game data...</div>;
  }

  if (!game) {
    return <div className="p-8 text-center text-red-500">Game not found</div>;
  }

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      <header className="bg-[var(--bg-secondary)] shadow-sm border-b border-[var(--border-color)]">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-[var(--text-primary)]">
              {game.away_team} @ {game.home_team}
            </h1>
            <p className="text-[var(--text-secondary)]">{[game.stadium, formatKickoff(game.kickoff)].filter(Boolean).join(' • ')}</p>
          </div>
          <div className="flex items-center gap-4">
            <Link to={`/season/${game.season}/week/${game.week}`} className="btn-ghost text-sm">
              ← Back
            </Link>
            <ThemeToggle />
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="space-y-6">
          <GameHeader game={game} />
          <Tabs
            tabs={[
              { id: 'stats', label: 'Stats' },
              { id: 'pbp', label: 'Play-by-Play' },
              { id: 'notes', label: 'Notes' },
            ]}
            defaultTab="stats"
          >
            <TabPanel id="stats">
              {stats && <GameStats stats={stats} homeTeam={game.home_team} awayTeam={game.away_team} />}
            </TabPanel>
            <TabPanel id="pbp">
              <DriveByPlayWrapper gameId={gameId} homeTeam={game.home_team} awayTeam={game.away_team} />
            </TabPanel>
            <TabPanel id="notes">
              <GameNotesPanel notes={notes} />
            </TabPanel>
          </Tabs>
        </div>
      </main>
    </div>
  );
}

function GameHeader({ game }: { game: Game }) {
  return (
    <div className="card p-6">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div className="flex items-center gap-6">
          <div className="flex items-center gap-4">
            <TeamLogo team={game.away_team} size="lg" />
            <div className="text-right">
              <p className="text-3xl font-bold text-nfl-dark dark:text-nfl-light">{game.away_score ?? '–'}</p>
              <p className="text-sm text-[var(--text-secondary)]">{game.away_team}</p>
            </div>
          </div>
          <div className="flex flex-col items-center px-4 border-x border-gray-200 dark:border-gray-700 my-4 md:my-0">
            <span className="text-sm text-[var(--text-secondary)]">{game.stadium}</span>
            <span className="text-xs text-[var(--text-secondary)]">{formatKickoff(game.kickoff)}</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="text-right">
              <p className="text-3xl font-bold text-nfl-light dark:text-nfl-light">{game.home_score ?? '–'}</p>
              <p className="text-sm text-[var(--text-secondary)]">{game.home_team}</p>
            </div>
            <TeamLogo team={game.home_team} size="lg" />
          </div>
        </div>
        <div className="text-right text-sm text-[var(--text-secondary)]">
          <p>Week {game.week} • {game.season} Season</p>
        </div>
      </div>
    </div>
  );
}

function GameStats({ stats, homeTeam, awayTeam }: { stats: TeamComparisonStats; homeTeam: string; awayTeam: string }) {
  return (
    <div className="space-y-6">
      <TeamComparisonCharts stats={stats} homeTeam={homeTeam} awayTeam={awayTeam} />

      <div className="grid gap-6 md:grid-cols-2">
        <StatSection
          title="Offense"
          home={stats.home_offense}
          away={stats.away_offense}
          homeTeam={homeTeam}
          awayTeam={awayTeam}
        />
        <StatSection
          title="Defense"
          defense
          home={stats.home_defense}
          away={stats.away_defense}
          homeTeam={homeTeam}
          awayTeam={awayTeam}
        />
      </div>
    </div>
  );
}

// A defense unit's numbers are what the opposing offense did against it, so lower is better,
// except takeaways.
function StatSection({
  title,
  home,
  away,
  homeTeam,
  awayTeam,
  defense = false,
}: {
  title: string;
  home: TeamUnitStats;
  away: TeamUnitStats;
  homeTeam: string;
  awayTeam: string;
  defense?: boolean;
}) {
  const higher = !defense;
  const teams = { homeTeam, awayTeam };
  const tdNote = (u: TeamUnitStats) => `${u.red_zone_tds}/${u.red_zone_trips}`;

  return (
    <div className="card p-4">
      <h4 className="font-semibold text-nfl-dark dark:text-nfl-light mb-4">{title}</h4>
      <div className="space-y-2">
        <StatRow label={defense ? 'EPA/Play Allowed' : 'EPA/Play'} home={home.epa_per_play} away={away.epa_per_play} format={fixed2} higherIsBetter={higher} {...teams} />
        <StatRow label={defense ? 'Success Rate Allowed' : 'Success Rate'} home={home.success_rate * 100} away={away.success_rate * 100} format={percent} higherIsBetter={higher} {...teams} />
        <StatRow label="Pass Rate" home={home.pass_rate * 100} away={away.pass_rate * 100} format={percent} higherIsBetter={higher} {...teams} />
        <StatRow label="Pass EPA" home={home.pass_epa} away={away.pass_epa} format={fixed2} higherIsBetter={higher} {...teams} />
        <StatRow label="Rush EPA" home={home.rush_epa} away={away.rush_epa} format={fixed2} higherIsBetter={higher} {...teams} />
        <StatRow
          label="Red Zone TD %"
          home={home.red_zone_pct * 100}
          away={away.red_zone_pct * 100}
          format={percent}
          homeNote={tdNote(home)}
          awayNote={tdNote(away)}
          higherIsBetter={higher}
          {...teams}
        />
        <StatRow label="3rd Down %" home={home.third_down_pct * 100} away={away.third_down_pct * 100} format={percent} higherIsBetter={higher} {...teams} />
        <StatRow label="4th Down %" home={home.fourth_down_pct * 100} away={away.fourth_down_pct * 100} format={percent} higherIsBetter={higher} {...teams} />
        {!defense && (
          <StatRow label="Time of Possession" home={home.time_of_possession_seconds} away={away.time_of_possession_seconds} format={mmss} {...teams} />
        )}
        <StatRow label={defense ? 'Explosive Plays Allowed' : 'Explosive Plays'} home={home.explosive_plays} away={away.explosive_plays} higherIsBetter={higher} {...teams} />
        <StatRow
          label={defense ? 'Takeaways' : 'Turnovers'}
          home={home.turnovers}
          away={away.turnovers}
          higherIsBetter={defense}
          {...teams}
        />
      </div>
    </div>
  );
}

function DriveByPlayWrapper({ gameId, homeTeam, awayTeam }: { gameId: string; homeTeam: string; awayTeam: string }) {
  const { data: drives, isLoading } = useDrives(gameId);

  if (isLoading) {
    return <div className="card p-8 text-center text-[var(--text-secondary)]">Loading drives...</div>;
  }

  if (!drives) {
    return <div className="card p-8 text-center text-[var(--text-secondary)]">No drive data available</div>;
  }

  return <DriveList data={drives} homeTeam={homeTeam} awayTeam={awayTeam} />;
}

function GameNotesPanel({ notes }: { notes: GameNotes | undefined }) {
  if (!notes) {
    return <div className="card p-8 text-center text-[var(--text-secondary)]">No notes data available</div>;
  }

  const fields: [string, string][] = [
    ['Stadium', notes.stadium],
    ['Surface', notes.surface],
    ['Roof', notes.roof],
    ['Indoor', notes.indoor ? 'Yes' : 'No'],
  ];

  return (
    <div className="card p-6">
      <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light mb-4">Game Information</h3>
      <div className="grid gap-4 md:grid-cols-2">
        {fields.map(([label, value]) => (
          <div key={label}>
            <p className="text-sm text-[var(--text-secondary)]">{label}</p>
            <p className="font-medium">{value || 'Not available'}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
