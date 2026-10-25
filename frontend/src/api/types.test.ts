import { describe, expect, it } from 'vitest';
import { mapTeamComparisonStats, type TeamUnitStatsAPI } from './types';

const unit = (overrides: Partial<TeamUnitStatsAPI> = {}): TeamUnitStatsAPI => ({
  plays: 0,
  epa_per_play: 0,
  success_rate: 0,
  pass_rate: 0,
  rush_epa: 0,
  pass_epa: 0,
  red_zone_trips: 0,
  red_zone_tds: 0,
  red_zone_pct: 0,
  third_down_pct: 0,
  fourth_down_pct: 0,
  time_of_possession_ns: 0,
  explosive_plays: 0,
  turnovers: 0,
  ...overrides,
});

describe('mapTeamComparisonStats', () => {
  it('formats time of possession nanoseconds as m:ss', () => {
    const mapped = mapTeamComparisonStats({
      home_team: 'GB',
      away_team: 'DAL',
      home_offense: unit({ time_of_possession_ns: 125_000_000_000 }),
      home_defense: unit(),
      away_offense: unit(),
      away_defense: unit(),
    });

    expect(mapped.home_offense.time_of_possession).toBe('2:05');
  });
});
