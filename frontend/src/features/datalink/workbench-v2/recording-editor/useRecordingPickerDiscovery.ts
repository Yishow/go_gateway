import { useEffect, useRef, useState } from 'react';
import type { GroupCandidate } from '../state/writeGroup/candidates';
export interface RecordingPickerDiscovery { candidates: readonly GroupCandidate[]; devices: readonly { id: string; name: string }[] }
/** Optional explicit-demand read projection; never hydrates or replaces a draft. */
export function useRecordingPickerDiscovery(scope: string, load?: () => Promise<RecordingPickerDiscovery>) {
  const owner = useRef({ scope }); if (owner.current.scope !== scope) owner.current = { scope };
  const current = useRef(load); current.current = load;
  const mounted = useRef(true), generation = useRef(0);
  const [result, setResult] = useState<{ scope: string; request: number; value?: RecordingPickerDiscovery; status: 'loading' | 'failed' | 'loaded' }>();
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; generation.current += 1; }; }, []);
  const refresh = async () => {
    if (!current.current) return;
    const captured = owner.current, request = ++generation.current;
    setResult({ scope: captured.scope, request, status: 'loading' });
    try {
      const value = await current.current();
      if (mounted.current && owner.current === captured && generation.current === request) setResult({ scope: captured.scope, request, status: 'loaded', value });
    } catch {
      if (mounted.current && owner.current === captured && generation.current === request) setResult({ scope: captured.scope, request, status: 'failed' });
    }
  };
  return { refresh, currentRequest: () => generation.current, result: result?.scope === scope ? result : undefined };
}
