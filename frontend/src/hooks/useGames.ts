import { useQuery } from '@tanstack/react-query';
import * as api from '../api/client';

export function useGames(season: number, week: number) {
  return useQuery({
    queryKey: ['games', season, week],
    queryFn: () => api.getGames(season, week),
    enabled: !!season && !!week,
  });
}

export function useWeeks(season: number) {
  return useQuery({
    queryKey: ['weeks', season],
    queryFn: () => api.getWeeks(season),
    enabled: !!season,
  });
}

export function useSeasons() {
  return useQuery({
    queryKey: ['seasons'],
    queryFn: api.getSeasons,
  });
}

export function useGame(gameId: string) {
  return useQuery({
    queryKey: ['game', gameId],
    queryFn: () => api.getGame(gameId),
    enabled: !!gameId,
  });
}

export function useGameInfo(gameId: string) {
  return useGame(gameId);
}

export function useGameNotes(gameId: string) {
  return useQuery({
    queryKey: ['gameNotes', gameId],
    queryFn: () => api.getGameNotes(gameId),
    enabled: !!gameId,
  });
}