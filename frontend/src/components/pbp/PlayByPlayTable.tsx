import { useState } from 'react';
import { QuarterAccordion } from './QuarterAccordion';
import type { PlayDetail } from '../../api/types';

export function PlayByPlayTable({ plays }: { plays: PlayDetail[] }) {
  const [quarterFilter, setQuarterFilter] = useState<number | 'all'>('all');
  const quarters = [1, 2, 3, 4].filter((q) => plays.some((p) => p.quarter === q));
  const filtered = quarterFilter === 'all' ? plays : plays.filter((p) => p.quarter === quarterFilter);

  return (
    <div>
      <div className="flex gap-2 mb-4 border-b pb-2">
        <button
          onClick={() => setQuarterFilter('all')}
          className={`px-3 py-1 rounded text-sm ${quarterFilter === 'all' ? 'bg-nfl-dark text-white' : 'bg-gray-100 dark:bg-gray-800 text-nfl-dark dark:text-nfl-light'}`}
        >
          All
        </button>
        {quarters.map((q) => (
          <button
            key={q}
            onClick={() => setQuarterFilter(q)}
            className={`px-3 py-1 rounded text-sm ${quarterFilter === q ? 'bg-nfl-medium text-white' : 'bg-gray-100 dark:bg-gray-800 text-nfl-dark dark:text-nfl-light'}`}
          >
            Q{q}
          </button>
        ))}
      </div>
      <div className="space-y-2">
        {quarterFilter === 'all' ? (
          quarters.map((q) => (
            <QuarterAccordion key={q} quarter={q} plays={plays.filter((p) => p.quarter === q)} />
          ))
        ) : (
          <QuarterAccordion quarter={quarterFilter} plays={filtered} />
        )}
      </div>
    </div>
  );
}