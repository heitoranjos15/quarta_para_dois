import ky from 'ky';
import type {
  Game,
  TeamComparisonStats,
  PlayDetail,
  GameNotes,
  APIResponse,
} from './types';

const api = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 30000,
  retry: 1,
});

export const getSeasons = () => api.get('api/seasons').json<number[]>();

export const getWeeks = (season: number) =>
  api.get(`api/seasons/${season}/weeks`).json<number[]>();

export const getGames = (season: number, week: number) =>
  api.get(`api/seasons/${season}/weeks/${week}/games`).json<Game[]>();

export const getGame = (gameId: string) =>
  api.get(`api/games/${gameId}`).json<Game>();

export const getGameStats = (gameId: string) =>
  api.get(`api/games/${gameId}/stats`).json<APIResponse<TeamComparisonStats>>();

export const getPlays = (
  gameId: string,
  params?: { quarter?: number; drive?: number }
) =>
  api.get(`api/games/${gameId}/plays`, { searchParams: params }).json<
    APIResponse<PlayDetail[]>
  >();

export const getGameNotes = (gameId: string) =>
  api.get(`api/games/${gameId}/notes`).json<APIResponse<GameNotes>>();