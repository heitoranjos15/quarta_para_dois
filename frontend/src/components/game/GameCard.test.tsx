import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it } from 'vitest';
import type { Game } from '../../api/types';
import { GameCard } from './GameCard';

const baseGame: Game = {
  game_id: '2026_06_DAL_GB',
  espn_id: '401873007',
  season: 2026,
  week: 6,
  game_type: 'REG',
  kickoff: '2026-10-19T00:20:00Z',
  status: 'final',
  home_team: 'GB',
  away_team: 'DAL',
  home_score: 23,
  away_score: 13,
  stadium: 'Lambeau Field',
  location: 'Green Bay, WI',
  indoor: false,
};

function renderCard(overrides: Partial<Game>) {
  return render(
    <MemoryRouter>
      <GameCard game={{ ...baseGame, ...overrides }} />
    </MemoryRouter>,
  );
}

describe('GameCard', () => {
  it('highlights only the winner of a decided game', () => {
    const { container } = renderCard({});

    const winners = container.querySelectorAll('.team-score-winner');
    expect(winners).toHaveLength(1);
    expect(winners[0].textContent).toContain('23');
  });

  it('highlights nobody when the game is tied', () => {
    const { container } = renderCard({ home_score: 20, away_score: 20 });

    expect(container.querySelectorAll('.team-score-winner')).toHaveLength(0);
  });

  it('highlights nobody and shows a dash before the game has scores', () => {
    const { container, getAllByText } = renderCard({
      status: 'scheduled',
      home_score: null,
      away_score: null,
    });

    expect(container.querySelectorAll('.team-score-winner')).toHaveLength(0);
    expect(getAllByText('–')).toHaveLength(2);
  });
});
