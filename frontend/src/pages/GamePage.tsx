import { useParams } from 'react-router-dom';
import { useGame } from '../hooks/useGames';
import { useGameStats } from '../hooks/useGameStats';
import { usePlays } from '../hooks/usePlays';
import type { PlayDetail, TeamUnitStats } from '../api/types';

export function GamePage() {
  const params = useParams();
  const gameId = params.gameId || '';

  const { data: game, isLoading: gameLoading } = useGame(gameId);
  const { data: stats, isLoading: statsLoading } = useGameStats(gameId);
  const { data: plays, isLoading: playsLoading } = usePlays(gameId);

  if (gameLoading || statsLoading || playsLoading) {
    return <div className="p-8 text-center">Loading game data...</div>;
  }

  if (!game) {
    return <div className="p-8 text-center text-red-500">Game not found</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <h1 className="text-2xl font-bold text-gray-900">
            {game.away_team} @ {game.home_team}
          </h1>
          <p className="text-gray-500">{game.venue} • {new Date(game.game_date).toLocaleDateString()}</p>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid gap-6 lg:grid-cols-3">
          <div className="lg:col-span-2 space-y-6">
            <section className="bg-white rounded-lg border border-gray-200 p-6">
              <h2 className="text-lg font-semibold mb-4">Game Stats</h2>
              {stats && (
                <div className="grid gap-4 md:grid-cols-2">
                  <StatCard title="Home Offense" stats={stats.home_offense} />
                  <StatCard title="Away Offense" stats={stats.away_offense} />
                  <StatCard title="Home Defense" stats={stats.home_defense} />
                  <StatCard title="Away Defense" stats={stats.away_defense} />
                </div>
              )}
            </section>

            <section className="bg-white rounded-lg border border-gray-200 p-6">
              <h2 className="text-lg font-semibold mb-4">Play by Play</h2>
              <div className="space-y-2 max-h-96 overflow-y-auto">
                {plays?.map((play: PlayDetail) => (
                  <div key={play.play_id} className="border-b py-2 last:border-0">
                    <div className="flex gap-4 text-sm">
                      <span className="font-mono w-20">{play.quarter}Q</span>
                      <span className="w-24">{play.down} & {play.distance}</span>
                      <span className="flex-1">{play.description}</span>
                      <span className="font-mono text-right w-16">{play.epa.toFixed(2)}</span>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          </div>

          <aside className="space-y-6">
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <h3 className="font-semibold mb-4">Score</h3>
              <div className="flex justify-between text-3xl font-bold">
                <span>{game.away_team}: {game.away_score}</span>
                <span>{game.home_team}: {game.home_score}</span>
              </div>
            </div>
          </aside>
        </div>
      </main>
    </div>
  );
}

function StatCard({ title, stats }: { title: string; stats: TeamUnitStats }) {
  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <h4 className="font-medium text-gray-700 mb-3">{title}</h4>
      <dl className="space-y-2 text-sm">
        <div className="flex justify-between"><dt>EPA/Play</dt><dd className="font-mono">{stats.epa_per_play.toFixed(2)}</dd></div>
        <div className="flex justify-between"><dt>Success Rate</dt><dd className="font-mono">{(stats.success_rate * 100).toFixed(1)}%</dd></div>
        <div className="flex justify-between"><dt>Pass Rate</dt><dd className="font-mono">{(stats.pass_rate * 100).toFixed(1)}%</dd></div>
        <div className="flex justify-between"><dt>Explosive Plays</dt><dd className="font-mono">{stats.explosive_plays}</dd></div>
        <div className="flex justify-between"><dt>Turnovers</dt><dd className="font-mono">{stats.turnovers}</dd></div>
      </dl>
    </div>
  );
}