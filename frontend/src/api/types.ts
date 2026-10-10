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

export interface GameNotes {
  weather: WeatherInfo;
  injuries: InjuryInfo[];
  officials: OfficialInfo[];
  game_info: GameInfo;
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