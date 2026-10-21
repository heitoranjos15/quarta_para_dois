import { useState } from 'react';
import { Link } from 'react-router-dom';
import { clsx } from 'clsx';
import type { Game } from '../../api/types';
import { TeamLogo } from '../ui/TeamLogo';
import { ChevronIcon } from '../ui/ChevronIcon';
import { Button } from '../ui/Button';

type Side = 'winner' | 'loser' | undefined;

// Only a finished game with different scores has a winner; ties and unplayed games have none.
function sides(game: Game): { home: Side; away: Side } {
  const { home_score: home, away_score: away } = game;
  if (game.status !== 'final' || home === null || away === null || home === away) {
    return { home: undefined, away: undefined };
  }
  return home > away ? { home: 'winner', away: 'loser' } : { home: 'loser', away: 'winner' };
}

function statusLine(game: Game): string {
  if (game.status === 'final') return 'Final';
  if (game.status === 'in_progress') return 'Live';
  return game.kickoff ? new Date(game.kickoff).toLocaleString() : '';
}

function scoreClass(side: Side) {
  return clsx(side === 'winner' && 'team-score-winner', side === 'loser' && 'team-score-loser');
}

export function GameCard({ game }: { game: Game }) {
  const [expanded, setExpanded] = useState(false);
  const result = sides(game);
  const status = statusLine(game);

  return (
    <article data-testid="game-card" data-game-id={game.game_id} className="game-card">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full p-4 flex items-center justify-between text-left"
        aria-expanded={expanded}
      >
        <div className="flex items-center gap-4 flex-1 min-w-0">
          <TeamLogo team={game.away_team} size="lg" />
          <div className={clsx('text-right min-w-[60px]', scoreClass(result.away))}>
            <p className="font-bold text-lg">{game.away_score ?? '–'}</p>
            <p className="text-sm text-[var(--text-muted)] truncate">{game.away_team}</p>
          </div>
          <span className="text-xl font-bold text-[var(--text-muted)] px-4 flex-shrink-0">×</span>
          <div className={clsx('text-left min-w-[60px]', scoreClass(result.home))}>
            <p className="font-bold text-lg">{game.home_score ?? '–'}</p>
            <p className="text-sm text-[var(--text-muted)] truncate">{game.home_team}</p>
          </div>
          <TeamLogo team={game.home_team} size="lg" />
          {status && <span className="text-xs text-[var(--text-muted)] ml-auto pr-2">{status}</span>}
        </div>
        <ChevronIcon className={clsx(expanded && 'rotate-180')} />
      </button>
      {expanded && (
        <div className="game-card-expanded px-4 py-3">
          <Link to={`/game/${game.game_id}`}>
            <Button variant="primary" className="w-full">
              View Game Details
            </Button>
          </Link>
        </div>
      )}
    </article>
  );
}
