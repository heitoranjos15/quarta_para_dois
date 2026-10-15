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
import type { TeamComparisonStats, TeamUnitStats, GameNotes } from '../api/types';

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
            <p className="text-[var(--text-secondary)]">{game.venue} • {new Date(game.game_date).toLocaleDateString()}</p>
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

function GameHeader({ game }: { game: { home_team: string; away_team: string; home_score: number; away_score: number; venue: string; game_date: string; season: number; week: number } }) {
  return (
    <div className="card p-6">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div className="flex items-center gap-6">
          <div className="flex items-center gap-4">
            <TeamLogo team={game.away_team} size="lg" />
            <div className="text-right">
              <p className="text-3xl font-bold text-nfl-dark dark:text-nfl-light">{game.away_score}</p>
              <p className="text-sm text-[var(--text-secondary)]">{game.away_team}</p>
            </div>
          </div>
          <div className="flex flex-col items-center px-4 border-x border-gray-200 dark:border-gray-700 my-4 md:my-0">
            <span className="text-sm text-[var(--text-secondary)]">{game.venue}</span>
            <span className="text-xs text-[var(--text-secondary)]">{new Date(game.game_date).toLocaleDateString()}</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="text-right">
              <p className="text-3xl font-bold text-nfl-light dark:text-nfl-light">{game.home_score}</p>
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
        <StatSection title={`${homeTeam} Offense`} stats={stats.home_offense} homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatSection title={`${awayTeam} Offense`} stats={stats.away_offense} homeTeam={awayTeam} awayTeam={homeTeam} />
        <StatSection title={`${homeTeam} Defense`} stats={stats.home_defense} homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatSection title={`${awayTeam} Defense`} stats={stats.away_defense} homeTeam={awayTeam} awayTeam={homeTeam} />
      </div>
    </div>
  );
}

function StatSection({ title, stats, homeTeam, awayTeam }: { title: string; stats: TeamUnitStats; homeTeam: string; awayTeam: string }) {
  return (
    <div className="card p-4">
      <h4 className="font-semibold text-nfl-dark dark:text-nfl-light mb-4">{title}</h4>
      <div className="space-y-2">
        <StatRow label="EPA/Play" home={stats.epa_per_play} away={0} unit="" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Success Rate" home={stats.success_rate * 100} away={0} unit="%" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Pass Rate" home={stats.pass_rate * 100} away={0} unit="%" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Pass EPA" home={stats.pass_epa} away={0} unit="" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Rush EPA" home={stats.rush_epa} away={0} unit="" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Red Zone %" home={stats.red_zone_pct * 100} away={0} unit="%" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="3rd Down %" home={stats.third_down_pct * 100} away={0} unit="%" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="4th Down %" home={stats.fourth_down_pct * 100} away={0} unit="%" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Time of Possession" home={0} away={0} unit="" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Explosive Plays" home={stats.explosive_plays} away={0} unit="" homeTeam={homeTeam} awayTeam={awayTeam} />
        <StatRow label="Turnovers" home={stats.turnovers} away={0} unit="" higherIsBetter={false} homeTeam={homeTeam} awayTeam={awayTeam} />
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

  return (
    <div className="space-y-6">
      <div className="card p-6">
        <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light mb-4">Game Information</h3>
        <div className="grid gap-4 md:grid-cols-2">
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Stadium</p>
            <p className="font-medium">{notes.game_info.stadium}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Surface</p>
            <p className="font-medium">{notes.game_info.surface}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Roof</p>
            <p className="font-medium">{notes.game_info.roof}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Attendance</p>
            <p className="font-medium">{notes.game_info.attendance.toLocaleString()}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Duration</p>
            <p className="font-medium">{notes.game_info.duration}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Referee</p>
            <p className="font-medium">{notes.officials.find((o) => o.role === 'Referee')?.name || '—'}</p>
          </div>
        </div>
      </div>

      <div className="card p-6">
        <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light mb-4">Weather</h3>
        <div className="grid gap-4 md:grid-cols-3">
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Temperature</p>
            <p className="font-medium">{notes.weather.temperature}°F</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Conditions</p>
            <p className="font-medium">{notes.weather.conditions}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Wind</p>
            <p className="font-medium">{notes.weather.wind_speed} mph {notes.weather.wind_direction}</p>
          </div>
          <div>
            <p className="text-sm text-[var(--text-secondary)]">Humidity</p>
            <p className="font-medium">{notes.weather.humidity}%</p>
          </div>
        </div>
      </div>

      <div className="card p-6">
        <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light mb-4">Injuries</h3>
        {notes.injuries.length === 0 ? (
          <p className="text-[var(--text-secondary)]">No injuries reported</p>
        ) : (
          <div className="divide-y divide-gray-100 dark:divide-gray-800">
            {notes.injuries.map((injury, i) => (
              <div key={i} className="py-3 flex flex-col md:flex-row md:items-center md:justify-between gap-2">
                <div className="flex items-center gap-3">
                  <TeamLogo team={injury.team} size="sm" />
                  <div>
                    <p className="font-medium">{injury.player}</p>
                    <p className="text-sm text-[var(--text-secondary)]">{injury.body_part} • {injury.status}</p>
                  </div>
                </div>
                <span className="text-sm text-[var(--text-secondary)]">{injury.practice}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}