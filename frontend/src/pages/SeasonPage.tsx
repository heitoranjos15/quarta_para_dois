import { useParams, Link } from 'react-router-dom';
import { useWeeks } from '../hooks/useGames';
import { Loading } from '../components/ui/Loading';

export function SeasonPage() {
  const params = useParams();
  const season = params.season ? parseInt(params.season, 10) : new Date().getFullYear();
  const { data: weeks, isLoading } = useWeeks(season);

  if (isLoading) return <Loading size="lg" />;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <h1 className="text-2xl font-bold text-gray-900">Quarta Para Dois - NFL Analytics</h1>
          <p className="text-gray-500 mt-1">Season {season}</p>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
          {weeks?.map((week) => (
            <Link
              key={week}
              to={`/season/${season}/week/${week}`}
              className="bg-white rounded-lg border border-gray-200 p-6 text-center hover:shadow-md hover:border-blue-300 transition-all"
            >
              <span className="text-3xl font-bold text-gray-900">Week {week}</span>
            </Link>
          ))}
          {(!weeks || weeks.length === 0) && (
            <div className="col-span-full text-center text-gray-500 py-12">
              No weeks found for this season
            </div>
          )}
        </div>
      </main>
    </div>
  );
}