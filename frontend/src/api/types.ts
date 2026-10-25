// '' when only games.csv is available and the game is past kickoff without a score.
export type GameStatus = 'scheduled' | 'in_progress' | 'final' | '';

export interface Game {
  game_id: string;
  espn_id: string;
  season: number;
  week: number;
  game_type: string;
  kickoff: string; // RFC3339 UTC, "" if unknown
  status: GameStatus;
  home_team: string;
  away_team: string;
  home_score: number | null;
  away_score: number | null;
  stadium: string;
  location: string;
  indoor: boolean;
  roof?: string;
  surface?: string;
}

export interface TeamComparisonStats {
  home_team: string;
  away_team: string;
  home_offense: TeamUnitStats;
  home_defense: TeamUnitStats;
  away_offense: TeamUnitStats;
  away_defense: TeamUnitStats;
}

export interface TeamUnitStats {
  plays: number;
  epa_per_play: number;
  success_rate: number;
  pass_rate: number;
  rush_epa: number;
  pass_epa: number;
  red_zone_trips: number;
  red_zone_tds: number;
  red_zone_pct: number;
  third_down_pct: number;
  fourth_down_pct: number;
  time_of_possession: string;
  time_of_possession_seconds: number;
  explosive_plays: number;
  turnovers: number;
}

export interface TeamComparisonStatsAPI {
  home_team: string;
  away_team: string;
  home_offense: TeamUnitStatsAPI;
  home_defense: TeamUnitStatsAPI;
  away_offense: TeamUnitStatsAPI;
  away_defense: TeamUnitStatsAPI;
}

export interface TeamUnitStatsAPI extends Omit<TeamUnitStats, 'time_of_possession' | 'time_of_possession_seconds'> {
  time_of_possession_ns: number;
}

export function mapTeamComparisonStats(api: TeamComparisonStatsAPI): TeamComparisonStats {
  return {
    home_team: api.home_team,
    away_team: api.away_team,
    home_offense: mapTeamUnitStats(api.home_offense),
    home_defense: mapTeamUnitStats(api.home_defense),
    away_offense: mapTeamUnitStats(api.away_offense),
    away_defense: mapTeamUnitStats(api.away_defense),
  };
}

function mapTeamUnitStats({ time_of_possession_ns, ...api }: TeamUnitStatsAPI): TeamUnitStats {
  return {
    ...api,
    time_of_possession: formatTimeOfPossession(time_of_possession_ns),
    time_of_possession_seconds: Math.floor(time_of_possession_ns / 1_000_000_000),
  };
}

function formatTimeOfPossession(nanos: number): string {
  const totalSeconds = Math.floor(nanos / 1_000_000_000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
}

export interface PlayDetail {
  play_id: number;
  quarter: number;
  down: number;
  distance: number;
  yardline: number;
  play_type: string;
  description: string;
  epa: number | null;
  passer: string;
  receiver: string;
  rusher: string;
  tackler: string;
  is_scoring_play: boolean;
  drive_id: number;
  drive_play_count: number;
  drive_yards: number;
  drive_result: string;
  drive_top: string;
}

export interface PlayDetailAPI extends Omit<PlayDetail, 'drive_top'> {
  drive_top_ns: number;
}

export function mapPlayDetail({ drive_top_ns, ...api }: PlayDetailAPI): PlayDetail {
  return {
    ...api,
    play_type: api.play_type.toLowerCase(),
    drive_top: formatTimeOfPossession(drive_top_ns),
  };
}

export interface PlayByPlayResponseAPI {
  quarter: number;
  plays: PlayDetailAPI[];
  drives: DriveAPI[] | null;
}

export interface DriveAPI {
  drive_id: number;
  start_yardline: number;
  end_yardline: number;
  yards: number;
  play_count: number;
  result: string;
  top_ns: number;
}

// Only what the API actually provides; Part 3 adds weather, injuries and officials.
export interface GameNotes {
  stadium: string;
  roof: string;
  surface: string;
  indoor: boolean;
}

export interface DriveDetail {
  drive_id: number;
  team: string;
  start_yardline: number;
  end_yardline: number;
  plays_count: number;
  yards: number;
  time_of_possession: string;
  result: string;
  points: number;
  end_play_description: string;
  end_play_type: string;
  plays: PlayDetail[];
}

export interface TeamDriveTotal {
  drives: number;
  points: number;
  time_of_possession: string;
  yards: number;
}

export interface QuarterDriveSummary {
  drives: DriveDetail[];
  total: {
    home: TeamDriveTotal;
    away: TeamDriveTotal;
  };
}

export interface DriveSummaryResponse {
  quarters: Record<number, QuarterDriveSummary>; // overtime is quarter 5
}

interface DriveDetailAPI extends Omit<DriveDetail, 'plays'> {
  plays: PlayDetailAPI[] | null;
}

interface QuarterDriveSummaryAPI {
  drives: DriveDetailAPI[] | null;
  total: QuarterDriveSummary['total'] | null;
}

export interface DriveSummaryResponseAPI {
  quarters: Record<number, QuarterDriveSummaryAPI> | null;
}

function mapQuarterDriveSummary(q: QuarterDriveSummaryAPI): QuarterDriveSummary {
  return {
    drives: (q.drives ?? []).map((d) => ({
      ...d,
      plays: (d.plays ?? []).map(mapPlayDetail),
    })),
    total: q.total ?? {
      home: { drives: 0, points: 0, time_of_possession: '0:00', yards: 0 },
      away: { drives: 0, points: 0, time_of_possession: '0:00', yards: 0 },
    },
  };
}

export function mapDriveSummary(api: DriveSummaryResponseAPI): DriveSummaryResponse {
  const quarters: Record<number, QuarterDriveSummary> = {};
  Object.entries(api.quarters ?? {}).forEach(([q, summary]) => {
    quarters[Number(q)] = mapQuarterDriveSummary(summary);
  });
  return { quarters };
}

export interface APIResponse<T> {
  data: T;
  meta: { cached: boolean; season: number };
}
