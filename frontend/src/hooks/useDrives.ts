import { useQuery } from '@tanstack/react-query';
import * as api from '../api/client';
import type { DriveSummaryResponse } from '../api/types';

export function useDrives(gameId: string) {
  return useQuery<DriveSummaryResponse>({
    queryKey: ['drives', gameId],
    queryFn: () => api.getDrives(gameId),
    enabled: !!gameId,
  });
}