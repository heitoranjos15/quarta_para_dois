import { describe, expect, it } from 'vitest';
import { mapTeamComparisonStats, type TeamUnitStatsAPI } from './types';

const unit = (overrides: Partial<TeamUnitStatsAPI> = {}): TeamUnitStatsAPI => ({
  Plays: 0,
  EPA_per_play: 0,
  SuccessRate: 0,
  PassRate: 0,
  RushEPA: 0,
  PassEPA: 0,
  RedZonePct: 0,
  ThirdDownPct: 0,
  FourthDownPct: 0,
  TimeOfPossession: 0,
  ExplosivePlays: 0,
  Turnovers: 0,
  ...overrides,
});

describe('mapTeamComparisonStats', () => {
  it('formats time of possession nanoseconds as m:ss', () => {
    const mapped = mapTeamComparisonStats({
      HomeTeam: 'GB',
      AwayTeam: 'DAL',
      HomeOffense: unit({ TimeOfPossession: 125_000_000_000 }),
      HomeDefense: unit(),
      AwayOffense: unit(),
      AwayDefense: unit(),
    });

    expect(mapped.home_offense.time_of_possession).toBe('2:05');
  });
});
