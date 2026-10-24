import type { Mapping, Point, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { WriteGroup, WriteGroupProfileKind } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from './writeGroupState';

/** One saved device with six float Tags: pt-0..pt-5 / tag-0..tag-5. */
export function batchPanelState(): WorkbenchV2State {
  const points: Point[] = [];
  const mappings: Record<string, Mapping> = {};
  for (let i = 0; i < 6; i++) {
    const pid = `rule-1-p-${i}`;
    const name = `sensor.${i}`;
    points.push({
      id: pid, device_id: 'dev-1', rule_id: 'rule-1', rule_name: 'Rule 1', name, address: `4000${i + 1}`,
      data_type: 'float64', function: 'holding_register', width: 1, enabled: true, skipped: false,
      _rule_scale: 1, _rule_offset: 0,
    });
    mappings[pid] = {
      point_id: pid, tag_key: name, display_name: name, unit: '', target_type: 'float64', scale: 1, offset: 0,
      enabled: true, persisted: true, save_state: 'saved', persisted_point_id: `pt-${i}`, tag_id: `tag-${i}`, device_id: 'dev-1',
    };
  }
  return { ...savedGroupState(), points, mappings };
}

/** A saved managed long group; cumulative intervals are minutes, instantaneous seconds. */
export function longProfileGroup(input: {
  id: string; circuitId: string; circuitName: string; interval: number; points: [number, number];
  kind?: WriteGroupProfileKind;
}): WriteGroup {
  const kind = input.kind ?? 'instantaneous';
  const seconds = kind === 'cumulative' ? input.interval * 60 : input.interval;
  return savedGroup({
    id: input.id,
    revision: `rev-${input.id}-1`,
    name: `${input.circuitName} (${kind})`,
    status: 'ready',
    layout_version: 'managed_long_v1',
    recording_profile: { kind, interval: input.interval, unit: kind === 'cumulative' ? 'minute' : 'second' },
    storage_mode: 'plain',
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_Safe_ws_1_long' },
    members: input.points.map((i) => ({
      device_id: 'dev-1', point_id: `pt-${i}`, tag_id: `tag-${i}`, source_revision: `s${i}`, mapping_revision: `m${i}`,
      target_column: '', required: false, circuit_id: input.circuitId, circuit_name: input.circuitName,
    })),
    row_policy: { interval_seconds: seconds, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
  });
}

/** An axios-shaped rejection that carries an HTTP status and an optional typed backend code. */
export function httpFailure(status: number, code?: string, message = `Request failed with status code ${status}`): Error {
  return Object.assign(new Error(message), { response: { status, data: code ? { error: { code } } : {} } });
}
