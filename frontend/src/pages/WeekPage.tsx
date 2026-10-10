import { useParams, Link } from 'react-router-dom';
import { useWeeks, useGames } from '../hooks/useGames';
import { GameCard } from '../components/game/GameCard';
import { Loading } from '../components/ui/Loading';
import { Button } from '../components/ui/Button';

export function WeekPage() {
  const params = useParams();
  const season = params.season ? parseInt(params.season, 10) : new Date().getFullYear();
  const week = params.week ? parseInt(params.week, 10) : 1;

  const { data: weeks, isLoading: weeksLoading } = useWeeks(season);
  const { data: games, isLoading: gamesLoading, error } = useGames(season, week);

  if (weeksLoading) return <Loading size="lg" />;
  if (gamesLoading) return <Loading size="lg" />;
  if (error) return <div className="p-8 text-center text-red-500">Error loading games</div>;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Week {week}</h1>
              <p className="text-gray-500">Season {season}</p>
            </div>
            <Link to={`/season/${season}`}>
              <Button variant="ghost" size="sm">← Back to Season</Button>
            </Link>
          </div>
        </div>
      </header>

      <nav className="max-w-7xl mx-auto px-4 py-4 bg-white border-b" aria-label="Week navigation">
        <div className="flex gap-2 overflow-x-auto pb-2">
          {weeks?.map((w) => (
            <Link
              key={w}
              to={`/season/${season}/week/${w}`}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${
                w === week
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Week {w}
            </Link>
          ))}
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="space-y-4">
          {games?.map((game) => (
            <GameCard key={game.game_id} game={game} />
          ))}
          {(!games || games.length === 0) && (
            <div className="text-center text-gray-500 py-12">
              No games found for this week
            </div>
          )}
        </div>
      </main>
    </div>
  );
}