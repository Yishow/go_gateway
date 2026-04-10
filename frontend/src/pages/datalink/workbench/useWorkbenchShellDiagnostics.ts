import { useCallback, useMemo, useState } from 'react';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { useWorkbench } from './WorkbenchProvider';

type ShellDiagnosticsState = 'idle' | 'running' | 'success' | 'error';

type RefetchableQuery = {
  refetch?: () => Promise<unknown>;
};

function safeRefetch(query: RefetchableQuery): Promise<unknown> {
  return query.refetch ? query.refetch() : Promise.resolve();
}

export function useWorkbenchShellDiagnostics() {
  const { selectedDeviceId } = useWorkbench();
  const devicesQuery = useDevicesQuery();
  const pointsQuery = usePointsQuery(selectedDeviceId ? { device_id: selectedDeviceId } : undefined);
  const tagsQuery = useTagsQuery();
  const mappingsQuery = useMappingsQuery();
  const [state, setState] = useState<ShellDiagnosticsState>('idle');

  const refreshDiagnostics = useCallback(async () => {
    setState('running');

    try {
      await Promise.all([
        safeRefetch(devicesQuery),
        safeRefetch(pointsQuery),
        safeRefetch(tagsQuery),
        safeRefetch(mappingsQuery),
      ]);
      setState('success');
    } catch {
      setState('error');
    }
  }, [devicesQuery, mappingsQuery, pointsQuery, tagsQuery]);

  return useMemo(
    () => ({
      refreshDiagnostics,
      statusKey: `workbench.shell.refresh.${state}`,
      actionKey:
        state === 'error'
          ? 'workbench.shell.refresh.retry'
          : 'workbench.shell.refresh.action',
      isRefreshing: state === 'running',
    }),
    [refreshDiagnostics, state],
  );
}
