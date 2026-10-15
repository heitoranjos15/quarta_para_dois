import ky from 'ky';
import type {
  GameAPI,
  TeamComparisonStatsAPI,
  PlayByPlayResponseAPI,
  GameNotesAPI,
  DriveSummaryResponse,
  APIResponse,
} from './types';
import {
  mapGameAPItoGame,
  mapTeamComparisonStats,
  mapPlayDetail,
  mapGameNotes,
} from './types';

const api = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 30000,
  retry: 1,
});

const extractData = <T,>(response: APIResponse<T>) => response.data;

export const getSeasons = () =>
  api.get('api/seasons').json<APIResponse<number[]>>().then(extractData);

export const getWeeks = (season: number) =>
  api.get(`api/seasons/${season}/weeks`).json<APIResponse<number[]>>().then(extractData);

export const getGames = (season: number, week: number) =>
  api
    .get(`api/seasons/${season}/weeks/${week}/games`)
    .json<APIResponse<GameAPI[]>>()
    .then(extractData)
    .then((games) => games.map(mapGameAPItoGame));

export const getGame = (gameId: string) =>
  api
    .get(`api/games/${gameId}`)
    .json<APIResponse<GameAPI>>()
    .then(extractData)
    .then(mapGameAPItoGame);

export const getGameStats = (gameId: string) =>
  api
    .get(`api/games/${gameId}/stats`)
    .json<APIResponse<TeamComparisonStatsAPI>>()
    .then(extractData)
    .then(mapTeamComparisonStats);

export const getPlays = (
  gameId: string,
  params?: { quarter?: number; drive?: number }
) =>
  api
    .get(`api/games/${gameId}/plays`, { searchParams: params })
    .json<APIResponse<PlayByPlayResponseAPI | PlayByPlayResponseAPI[]>>()
    .then(extractData)
    .then((response) => {
      // Handle both single quarter (when quarter param provided) and all quarters (array)
      const quarters = Array.isArray(response) ? response : [response];
      return quarters.flatMap((q) => q.Plays.map(mapPlayDetail));
    });

export const getGameNotes = (gameId: string) =>
  api
    .get(`api/games/${gameId}/notes`)
    .json<APIResponse<GameNotesAPI>>()
    .then(extractData)
    .then(mapGameNotes);

export const getDrives = (gameId: string) =>
  api
    .get(`api/games/${gameId}/drives`)
    .json<APIResponse<DriveSummaryResponse>>()
    .then(extractData);