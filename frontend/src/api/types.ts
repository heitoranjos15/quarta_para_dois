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
  HomeTeam: string;
  AwayTeam: string;
  HomeOffense: TeamUnitStatsAPI;
  HomeDefense: TeamUnitStatsAPI;
  AwayOffense: TeamUnitStatsAPI;
  AwayDefense: TeamUnitStatsAPI;
}

export interface TeamUnitStatsAPI {
  Plays: number;
  EPA_per_play: number;
  SuccessRate: number;
  PassRate: number;
  RushEPA: number;
  PassEPA: number;
  RedZoneTrips: number;
  RedZoneTDs: number;
  RedZonePct: number;
  ThirdDownPct: number;
  FourthDownPct: number;
  TimeOfPossession: number;
  ExplosivePlays: number;
  Turnovers: number;
}

export function mapTeamComparisonStats(api: TeamComparisonStatsAPI): TeamComparisonStats {
  return {
    home_team: api.HomeTeam,
    away_team: api.AwayTeam,
    home_offense: mapTeamUnitStats(api.HomeOffense),
    home_defense: mapTeamUnitStats(api.HomeDefense),
    away_offense: mapTeamUnitStats(api.AwayOffense),
    away_defense: mapTeamUnitStats(api.AwayDefense),
  };
}

function mapTeamUnitStats(api: TeamUnitStatsAPI): TeamUnitStats {
  return {
    plays: api.Plays,
    epa_per_play: api.EPA_per_play,
    success_rate: api.SuccessRate,
    pass_rate: api.PassRate,
    rush_epa: api.RushEPA,
    pass_epa: api.PassEPA,
    red_zone_trips: api.RedZoneTrips,
    red_zone_tds: api.RedZoneTDs,
    red_zone_pct: api.RedZonePct,
    third_down_pct: api.ThirdDownPct,
    fourth_down_pct: api.FourthDownPct,
    time_of_possession: formatTimeOfPossession(api.TimeOfPossession),
    time_of_possession_seconds: Math.floor(api.TimeOfPossession / 1_000_000_000),
    explosive_plays: api.ExplosivePlays,
    turnovers: api.Turnovers,
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

export interface PlayDetailAPI {
  PlayID: number;
  Quarter: number;
  Down: number;
  Distance: number;
  Yardline: number;
  PlayType: string;
  Description: string;
  EPA: number | null;
  Passer: string;
  Receiver: string;
  Rusher: string;
  Tackler: string;
  IsScoringPlay: boolean;
  DriveID: number;
  DrivePlayCount: number;
  DriveYards: number;
  DriveResult: string;
  DriveTOP: string;
}

export function mapPlayDetail(api: PlayDetailAPI): PlayDetail {
  return {
    play_id: api.PlayID,
    quarter: api.Quarter,
    down: api.Down,
    distance: api.Distance,
    yardline: api.Yardline,
    play_type: api.PlayType.toLowerCase(),
    description: api.Description,
    epa: api.EPA,
    passer: api.Passer,
    receiver: api.Receiver,
    rusher: api.Rusher,
    tackler: api.Tackler,
    is_scoring_play: api.IsScoringPlay,
    drive_id: api.DriveID,
    drive_play_count: api.DrivePlayCount,
    drive_yards: api.DriveYards,
    drive_result: api.DriveResult,
    drive_top: api.DriveTOP,
  };
}

export interface PlayByPlayResponseAPI {
  Quarter: number;
  Plays: PlayDetailAPI[];
  Drives: DriveAPI[];
}

export interface DriveAPI {
  DriveID: number;
  StartYardline: number;
  EndYardline: number;
  Yards: number;
  PlayCount: number;
  Result: string;
  TOP: number;
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
