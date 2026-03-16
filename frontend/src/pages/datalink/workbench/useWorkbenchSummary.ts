import { useMemo } from 'react';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { useWorkbench } from './WorkbenchProvider';

export function useWorkbenchSummary() {
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
    const linkedTagCount = tags.filter((tag) => linkedTagIds.has(tag.id)).length;

    return {
      selectedDevice,
      pointCount: points.length,
      linkedTagCount,
      outputCandidateCount: linkedTagCount,
      sourceReady: points.length > 0,
      tagReady: linkedTagCount > 0,
      outputReady: linkedTagCount > 0,
    };
  }, [devices, mappings, points, selectedDeviceId, tags]);
}
