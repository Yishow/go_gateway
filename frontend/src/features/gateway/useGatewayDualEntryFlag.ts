import { useQuery } from '@tanstack/react-query';
import { fetchGatewayDualEntryEnabled } from './dualEntryFlag';

export function useGatewayDualEntryFlag() {
  return useQuery({
    queryKey: ['gateway', 'feature-flags', 'dual-entry'],
    queryFn: fetchGatewayDualEntryEnabled,
    staleTime: 30_000,
  });
}
