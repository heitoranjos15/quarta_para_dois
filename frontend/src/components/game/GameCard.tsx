import { useState } from 'react';
import { Link } from 'react-router-dom';
import { clsx } from 'clsx';
import type { Game } from '../../api/types';
import { TeamLogo } from '../ui/TeamLogo';
import { ChevronIcon } from '../ui/ChevronIcon';
import { Button } from '../ui/Button';

export function GameCard({ game }: { game: Game }) {
  const [expanded, setExpanded] = useState(false);
  const isHomeWinner = game.home_score > game.away_score;

  return (
    <article data-testid="game-card" data-game-id={game.game_id} className="game-card">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full p-4 flex items-center justify-between text-left"
        aria-expanded={expanded}
      >
        <div className="flex items-center gap-4 flex-1 min-w-0">
          <TeamLogo team={game.away_team} size="lg" />
          <div className={clsx('text-right min-w-[60px]', !isHomeWinner && 'team-score-winner', isHomeWinner && 'team-score-loser')}>
            <p className="font-bold text-lg">{game.away_score}</p>
            <p className="text-sm text-[var(--text-muted)] truncate">{game.away_team}</p>
          </div>
          <span className="text-xl font-bold text-[var(--text-muted)] px-4 flex-shrink-0">×</span>
          <div className={clsx('text-left min-w-[60px]', isHomeWinner && 'team-score-winner', !isHomeWinner && 'team-score-loser')}>
            <p className="font-bold text-lg">{game.home_score}</p>
            <p className="text-sm text-[var(--text-muted)] truncate">{game.home_team}</p>
          </div>
          <TeamLogo team={game.home_team} size="lg" />
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