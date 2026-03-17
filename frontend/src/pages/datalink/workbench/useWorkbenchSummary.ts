import { useMemo } from 'react';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { useWorkbench } from './WorkbenchProvider';
import {
  deviceStatusToReadiness,
  tagStatusToReadiness,
  type StepReadinessState,
} from './workbenchTypes';

export type WorkbenchSummary = {
  // Device context
  selectedDevice: ReturnType<typeof useDevicesQuery>['data'] extends
    | (infer D)[]
    | undefined
    ? D | null
    : never;

  // Counts
  pointCount: number;
  linkedTagCount: number;
  outputCandidateCount: number;

  // Legacy boolean flags (kept for backward-compat with existing consumers)
  sourceReady: boolean;
  tagReady: boolean;
  outputReady: boolean;

  // Rich per-step readiness (spec §8.2)
  deviceReadiness: StepReadinessState;
  sourceReadiness: StepReadinessState;
  tagReadiness: StepReadinessState;
  outputReadiness: StepReadinessState;
};

export function useWorkbenchSummary(): WorkbenchSummary {
  const { selectedDeviceId } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const { data: mappings = [] } = useMappingsQuery();
  const { data: tags = [] } = useTagsQuery();

  return useMemo(() => {
    const selectedDevice =
      devices.find((device) => device.id === selectedDeviceId) ?? null;
    const pointIds = new Set(points.map((point) => point.id));
    const linkedMappings = mappings.filter((mapping) => pointIds.has(mapping.point_id));
    const linkedTagIds = new Set(linkedMappings.map((mapping) => mapping.tag_id));
    const linkedTags = tags.filter((tag) => linkedTagIds.has(tag.id));
    const linkedTagCount = linkedTags.length;

    // --- Device readiness ---
    const deviceReadiness: StepReadinessState = selectedDevice
      ? { status: deviceStatusToReadiness(selectedDevice.status) }
      : { status: 'draft', reason: 'no-device-selected' };

    // --- Source readiness ---
    const sourceReady = points.length > 0;
    const sourceReadiness: StepReadinessState = !selectedDevice
      ? { status: 'draft', reason: 'no-device-selected' }
      : sourceReady
        ? { status: 'ready' }
        : { status: 'draft', reason: 'no-points' };

    // --- Tag readiness ---
    const tagReady = linkedTagCount > 0;
    const hasBlockedTags = linkedTags.some((tag) =>
      tagStatusToReadiness(tag.status, true) === 'blocked',
    );
    const tagReadiness: StepReadinessState = !selectedDevice
      ? { status: 'draft', reason: 'no-device-selected' }
      : !sourceReady
        ? { status: 'blocked', reason: 'no-source-points' }
        : tagReady
        ? hasBlockedTags
          ? { status: 'partial', reason: 'some-tags-blocked' }
          : { status: 'ready' }
        : { status: 'draft', reason: 'no-tags-linked' };

    // --- Output readiness ---
    // For now, output readiness mirrors tag readiness. When real output
    // mapping state (Local Modbus / Database) becomes available, this
    // will also incorporate per-target sync status.
    const outputReady = linkedTagCount > 0;
    const outputReadiness: StepReadinessState = !selectedDevice
      ? { status: 'draft', reason: 'no-device-selected' }
      : !sourceReady
        ? { status: 'blocked', reason: 'no-source-points' }
        : !tagReady
          ? { status: 'blocked', reason: 'no-tags-linked' }
          : { status: 'partial', reason: 'output-not-applied' };

    return {
      selectedDevice,
      pointCount: points.length,
      linkedTagCount,
      outputCandidateCount: linkedTagCount,
      sourceReady,
      tagReady,
      outputReady,
      deviceReadiness,
      sourceReadiness,
      tagReadiness,
      outputReadiness,
    };
  }, [devices, mappings, points, selectedDeviceId, tags]);
}
