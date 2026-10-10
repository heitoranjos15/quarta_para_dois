import { useQuery } from '@tanstack/react-query';
import * as api from '../api/client';
import type { TeamComparisonStats, APIResponse } from '../api/types';

export function useGameStats(gameId: string) {
  return useQuery({
    queryKey: ['gameStats', gameId],
    queryFn: () => api.getGameStats(gameId),
    enabled: !!gameId,
    select: (response: APIResponse<TeamComparisonStats>) => response.data,
  });
}