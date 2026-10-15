import { useParams, Link } from 'react-router-dom';
import { useWeeks } from '../hooks/useGames';
import { Loading } from '../components/ui/Loading';
import { ThemeToggle } from '../components/ui/ThemeToggle';

export function SeasonPage() {
  const params = useParams();
  const season = params.season ? parseInt(params.season, 10) : new Date().getFullYear();
  const { data: weeks, isLoading } = useWeeks(season);

  if (isLoading) return <Loading size="lg" />;

  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      <header className="bg-[var(--bg-secondary)] shadow-sm border-b border-[var(--border-color)]">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-[var(--text-primary)]">Quarta Para Dois - NFL Analytics</h1>
            <p className="text-[var(--text-secondary)] mt-1">Season {season}</p>
          </div>
          <ThemeToggle />
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
          {weeks?.map((week) => (
            <Link
              key={week}
              to={`/season/${season}/week/${week}`}
              className="week-card"
            >
              <span className="text-3xl font-bold text-[var(--text-primary)]">Week {week}</span>
            </Link>
          ))}
          {(!weeks || weeks.length === 0) && (
            <div className="col-span-full text-center text-[var(--text-muted)] py-12">
              No weeks found for this season
            </div>
          )}
        </div>
      </main>
    </div>
  );
}