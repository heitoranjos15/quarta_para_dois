export interface Game {
  game_id: string;
  season: number;
  week: number;
  home_team: string;
  away_team: string;
  home_score: number;
  away_score: number;
  game_date: string;
  venue: string;
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
  red_zone_pct: number;
  third_down_pct: number;
  fourth_down_pct: number;
  time_of_possession: string;
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
    red_zone_pct: api.RedZonePct,
    third_down_pct: api.ThirdDownPct,
    fourth_down_pct: api.FourthDownPct,
    time_of_possession: formatTimeOfPossession(api.TimeOfPossession),
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
  epa: number;
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
  EPA: number;
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

export interface GameNotes {
  weather: WeatherInfo;
  injuries: InjuryInfo[];
  officials: OfficialInfo[];
  game_info: GameInfo;
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
  quarters: Record<number, QuarterDriveSummary>;
  overtime?: QuarterDriveSummary;
}

export interface GameNotesAPI {
  weather: WeatherInfoAPI;
  vegas: VegasInfoAPI;
}

export interface WeatherInfoAPI {
  temperature: number;
  wind: number;
  condition: string;
  stadium: string;
  roof: string;
}

export interface VegasInfoAPI {
  spread: number;
  total: number;
  home_ml: number;
  away_ml: number;
}

export function mapGameNotes(api: GameNotesAPI): GameNotes {
  return {
    weather: {
      temperature: api.weather.temperature,
      wind_speed: api.weather.wind,
      wind_direction: '',
      conditions: api.weather.condition,
      humidity: 0,
    },
    injuries: [],
    officials: [],
    game_info: {
      game_id: '',
      stadium: api.weather.stadium,
      surface: '',
      roof: api.weather.roof,
      attendance: 0,
      duration: '',
    },
  };
}

export interface WeatherInfo {
  temperature: number;
  wind_speed: number;
  wind_direction: string;
  conditions: string;
  humidity: number;
}

export interface InjuryInfo {
  player: string;
  team: string;
  status: string;
  body_part: string;
  practice: string;
}

export interface OfficialInfo {
  name: string;
  role: string;
}

export interface GameInfo {
  game_id: string;
  stadium: string;
  surface: string;
  roof: string;
  attendance: number;
  duration: string;
}

export interface APIResponse<T> {
  data: T;
  meta: { cached: boolean; season: number };
}

export interface GameAPI {
  GameID: string;
  Season: number;
  Week: number;
  GameType: string;
  GameDate: string;
  Weekday: string;
  GameTime: string;
  HomeTeam: string;
  AwayTeam: string;
  HomeScore: number;
  AwayScore: number;
  Location: string;
  Stadium: string;
  Surface: string;
  Roof: string;
  Temp: number;
  Wind: number;
  HomeCoach: string;
  AwayCoach: string;
  Referee: string;
  StadiumId: string;
}

export function mapGameAPItoGame(api: GameAPI): Game {
  return {
    game_id: api.GameID,
    season: api.Season,
    week: api.Week,
    home_team: api.HomeTeam,
    away_team: api.AwayTeam,
    home_score: api.HomeScore,
    away_score: api.AwayScore,
    game_date: api.GameDate,
    venue: api.Stadium,
  };
}