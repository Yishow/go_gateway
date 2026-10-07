import type { RuntimeLogRecord } from '../../types/runtimeLogs';
export const BROWSER_LOG_MAX_RECORDS = 1000;
export const BROWSER_LOG_MAX_BYTES = 2 * 1024 * 1024;
const encoder = new TextEncoder();
export const runtimeLogBytes = (record: RuntimeLogRecord) => encoder.encode(JSON.stringify(record)).byteLength;
export const runtimeLogCursor = (record: RuntimeLogRecord) => `${record.instance_id}:${record.sequence}`;

export function appendRuntimeLogs(current: RuntimeLogRecord[], incoming: RuntimeLogRecord[]) {
  const records = current.slice();
  const ids = new Set(records.map(runtimeLogCursor));
  let bytes = records.reduce((sum, record) => sum + runtimeLogBytes(record), 0);
  let truncated = false;
  for (const record of incoming) {
    const id = runtimeLogCursor(record);
    if (ids.has(id)) continue;
    const size = runtimeLogBytes(record);
    if (size > BROWSER_LOG_MAX_BYTES) { truncated = true; continue; }
    records.push(record); ids.add(id); bytes += size;
    while (records.length > BROWSER_LOG_MAX_RECORDS || bytes > BROWSER_LOG_MAX_BYTES) {
      const removed = records.shift()!;
      bytes -= runtimeLogBytes(removed); ids.delete(runtimeLogCursor(removed)); truncated = true;
    }
  }
  return { records, bytes, truncated };
}
