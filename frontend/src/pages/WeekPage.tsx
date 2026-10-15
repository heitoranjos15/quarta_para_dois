import { useParams, Link, useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { useWeeks, useGames } from '../hooks/useGames';
import { GameCard } from '../components/game/GameCard';
import { Loading } from '../components/ui/Loading';
import { Button } from '../components/ui/Button';
import { ThemeToggle } from '../components/ui/ThemeToggle';
import { getCurrentNFLWeek, getCurrentNFLSeason } from '../utils/nfl';

export function WeekPage() {
  const params = useParams();
  const navigate = useNavigate();
  const season = params.season ? parseInt(params.season, 10) : getCurrentNFLSeason();
  const currentWeek = getCurrentNFLWeek(season);
  const week = params.week ? parseInt(params.week, 10) : currentWeek;

  useEffect(() => {
    if (!params.week) {
      navigate(`/season/${season}/week/${currentWeek}`, { replace: true });
    }
  }, [params.week, season, currentWeek, navigate]);

  const { data: weeks, isLoading: weeksLoading } = useWeeks(season);
  const { data: games, isLoading: gamesLoading, error } = useGames(season, week);

  if (weeksLoading) return <Loading size="lg" />;
  if (gamesLoading) return <Loading size="lg" />;
  if (error) return <div className="p-8 text-center text-red-500">Error loading games</div>;

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      <header className="bg-[var(--bg-secondary)] shadow-sm border-b border-[var(--border-color)]">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-[var(--text-primary)]">Week {week}</h1>
              <p className="text-[var(--text-secondary)]">Season {season}</p>
            </div>
            <div className="flex items-center gap-4">
              <Link to={`/season/${season}`}>
                <Button variant="ghost" size="sm">← Back to Season</Button>
              </Link>
              <ThemeToggle />
            </div>
          </div>
        </div>
      </header>

      <nav className="max-w-7xl mx-auto px-4 py-4 bg-[var(--bg-secondary)] border-b border-[var(--border-color)]" aria-label="Week navigation">
        <div className="flex gap-2 overflow-x-auto pb-2">
          {weeks?.map((w) => (
            <Link
              key={w}
              to={`/season/${season}/week/${w}`}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${
                w === week
                  ? 'bg-blue-600 text-white'
                  : 'bg-[var(--bg-tertiary)] text-[var(--text-secondary)] hover:bg-[var(--border-hover)]'
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
            <div className="text-center text-[var(--text-muted)] py-12">
              No games found for this week
            </div>
          )}
        </div>
      </main>
    </div>
  );
}