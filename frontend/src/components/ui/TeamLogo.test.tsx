import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { TeamLogo } from './TeamLogo';

describe('TeamLogo', () => {
  it('renders the full abbreviation so teams sharing a prefix look different', () => {
    const teams = ['NYG', 'NYJ', 'LAC', 'LA', 'LV'];
    render(
      <>
        {teams.map((t) => (
          <TeamLogo key={t} team={t} />
        ))}
      </>,
    );

    const texts = teams.map((t) => screen.getByLabelText(t).textContent);
    expect(texts).toEqual(teams);
    expect(new Set(texts).size).toBe(5);
  });

  it.each(['LA', 'WAS'])('uses team colors for %s, not the grey fallback', (team) => {
    render(<TeamLogo team={team} />);

    const primary = screen.getByLabelText(team).style.getPropertyValue('--team-primary');
    expect(primary).not.toBe('');
    expect(primary.toLowerCase()).not.toBe('#6b7280');
  });
});
