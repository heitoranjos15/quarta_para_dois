import { useQuery } from '@tanstack/react-query';
import * as api from '../api/client';

export function usePlays(gameId: string, quarter?: number, drive?: number) {
  return useQuery({
    queryKey: ['plays', gameId, quarter, drive],
    queryFn: () => api.getPlays(gameId, { quarter, drive }),
    enabled: !!gameId,
  });
}