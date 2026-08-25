import * as React from 'react';
import { useQueries, useQuery } from '@tanstack/react-query';
import { pointAPI } from '../../../../../services/datalink';
import type {
  PollResult,
  Point as PersistedPoint,
  RuntimeStreamConnectionState,
  RuntimeStreamRecovery,
  RuntimeValueEvent,
} from '../../../../../types/datalink';
import type { Mapping, Point } from '../../state/types';
import { pointKeys } from '../../../../../hooks/datalink/keys';

import {
  buildGroups,
  buildSubscriptions,
} from './step3LiveSubscription';
import { useStep3RuntimeStreams } from './useStep3RuntimeStreams';

export interface Step3LiveValuesState {
  connectionByDevice: Record<string, RuntimeStreamConnectionState>;
  liveValues: Record<string, RuntimeValueEvent>;
  rawValues: Record<string, unknown>;
  lastSuccessAtByDevice: Record<string, string>;
  recoveryByDevice: Record<string, RuntimeStreamRecovery>;
}

/**
 * 依賴陣列策略說明：
 * 本 hook 以「語意 key」（pointIdentityKey / mappingPersistedKey / snapshotPointDataKey /
 * subscriptionKey 等）取代物件參考作為 useMemo/useEffect 依賴。
 * 原因：points / mappings / pointQueries（useQueries 每次渲染回傳新陣列）等物件
 * 參考不穩定，若直接作為依賴會導致 subscription 重建與 SSE 重新連線。
 * 因此以下的 exhaustive-deps 警告為刻意例外，以 eslint-disable 標注。
 */
export function useStep3LiveValues(
  points: Point[],
  mappings: Record<string, Mapping>,
): Step3LiveValuesState {
  const pointIdentityKey = points
    .map((point) => `${point.id}:${point.device_id}:${point.address}`)
    .join('|');
  const mappingPersistedKey = Object.entries(mappings)
    .sort(([left], [right]) => left.localeCompare(right))
    .map(([pointId, mapping]) => `${pointId}:${mapping.persisted_point_id ?? ''}`)
    .join('|');
  const deviceIds = React.useMemo(
    () => Array.from(new Set(points.map((point) => point.device_id))).sort(),
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
    [pointIdentityKey],
  );
  const pointAddressToLocalPointId = React.useMemo(() => {
    const lookup: Record<string, string> = {};
    points.forEach((point) => {
      lookup[`${point.device_id}::${point.address}`] = point.id;
    });
    return lookup;
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
  }, [pointIdentityKey]);
  const mappingPersistedPointIdByLocalPointId = React.useMemo(() => {
    const lookup: Record<string, string> = {};
    Object.entries(mappings).forEach(([pointId, mapping]) => {
      if (mapping.persisted_point_id) {
        lookup[pointId] = mapping.persisted_point_id;
      }
    });
    return lookup;
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
  }, [mappingPersistedKey]);
  const pointQueries = useQueries({
    queries: deviceIds.map((deviceId) => ({
      queryKey: pointKeys.list({ device_id: deviceId }),
      queryFn: () => pointAPI.list({ device_id: deviceId }),
      enabled: deviceId.length > 0,
      retry: false,
      refetchOnWindowFocus: false,
    })),
  });
  const snapshotPointDataKey = React.useMemo(
    () =>
      pointQueries
        .map((query) =>
          (query.data ?? [])
            .map(
              (point: PersistedPoint) =>
                `${point.id}:${point.device_id}:${point.address}:${String(point.last_value ?? '')}:${point.last_read_at ?? ''}`,
            )
            .join(','),
        )
        .join('|'),
    [pointQueries],
  );
  const snapshotPersistedPointIdByLocalPointId = React.useMemo(() => {
    const lookup: Record<string, string> = {};
    pointQueries.forEach((query) => {
      (query.data ?? []).forEach((point: PersistedPoint) => {
        const localPointId = pointAddressToLocalPointId[`${point.device_id}::${point.address}`];
        if (localPointId) {
          lookup[localPointId] = point.id;
        }
      });
    });
    return lookup;
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
  }, [pointAddressToLocalPointId, snapshotPointDataKey]);
  const persistedPointIdByLocalPointId = React.useMemo(
    () => ({
      ...snapshotPersistedPointIdByLocalPointId,
      ...mappingPersistedPointIdByLocalPointId,
    }),
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
    [mappingPersistedKey, snapshotPointDataKey],
  );
  const subscriptionKey = React.useMemo(
    () =>
      points
        .map((point) => {
          const persistedPointId = persistedPointIdByLocalPointId[point.id] ?? point.id;
          return `${point.id}:${point.device_id}:${point.address}:${persistedPointId}`;
        })
        .join('|'),
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
    [pointIdentityKey, mappingPersistedKey, snapshotPointDataKey],
  );
  const subscriptions = React.useMemo(
    () => buildSubscriptions(points, persistedPointIdByLocalPointId),
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
    [subscriptionKey],
  );
  const groups = React.useMemo(() => buildGroups(subscriptions), [subscriptions]);
  const persistedToLocalPointId = React.useMemo(() => {
    const lookup: Record<string, string> = {};
    subscriptions.forEach((subscription) => {
      if (subscription.streamPointId !== subscription.localPointId) {
        lookup[subscription.streamPointId] = subscription.localPointId;
      }
    });
    return lookup;
  }, [subscriptions]);
  const pollPointIds = React.useMemo(
    () =>
      Array.from(
        new Set(
          subscriptions
            .map((subscription) => {
              const persistedPointId = subscription.streamPointId;
              return persistedPointId && persistedPointId !== subscription.localPointId
                ? persistedPointId
                : null;
            })
            .filter((pointId): pointId is string => Boolean(pointId)),
        ),
      ),
    [subscriptions],
  );
  const pollPointIdsKey = React.useMemo(() => pollPointIds.join(','), [pollPointIds]);
  const polledValuesQuery = useQuery({
    queryKey: ['workbench-v2', 'step3-point-poll', pollPointIdsKey],
    queryFn: () => pointAPI.batchPoll(pollPointIds),
    enabled: pollPointIds.length > 0,
    retry: false,
    refetchOnWindowFocus: false,
  });
  const {
    connectionByDevice,
    liveValues,
    lastSuccessAtByDevice,
    recoveryByDevice,
  } = useStep3RuntimeStreams(groups, persistedToLocalPointId, pointAddressToLocalPointId);
  const snapshotRawValues = React.useMemo(() => {
    const values: Record<string, unknown> = {};
    pointQueries.forEach((query) => {
      (query.data ?? []).forEach((point: PersistedPoint) => {
        const localPointId =
          persistedToLocalPointId[point.id] ??
          pointAddressToLocalPointId[`${point.device_id}::${point.address}`];
        if (!localPointId) {
          return;
        }
        if (point.last_value !== null && point.last_value !== undefined) {
          values[localPointId] = point.last_value;
        }
      });
    });
    return values;
    // eslint-disable-next-line react-hooks/exhaustive-deps -- 以語意 key 作為依賴（見檔頭說明）
  }, [persistedToLocalPointId, pointAddressToLocalPointId, snapshotPointDataKey]);
  const polledRawValues = React.useMemo(() => {
    const values: Record<string, unknown> = {};
    (polledValuesQuery.data ?? []).forEach((result: PollResult) => {
      const localPointId = persistedToLocalPointId[result.point_id];
      if (!localPointId || result.value === null || result.value === undefined || result.error) {
        return;
      }
      values[localPointId] = result.value;
    });
    return values;
  }, [persistedToLocalPointId, polledValuesQuery.data]);
  const rawValues = React.useMemo(() => {
    const values: Record<string, unknown> = {
      ...snapshotRawValues,
      ...polledRawValues,
    };
    Object.entries(liveValues).forEach(([pointId, event]) => {
      values[pointId] = event.raw_value;
    });
    return values;
  }, [liveValues, polledRawValues, snapshotRawValues]);

  return {
    connectionByDevice,
    liveValues,
    rawValues,
    lastSuccessAtByDevice,
    recoveryByDevice,
  };
}
