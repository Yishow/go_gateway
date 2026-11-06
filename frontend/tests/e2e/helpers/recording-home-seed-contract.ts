export interface SeedDevice { id: string; name: string; protocol: string; status: string; running: boolean; availability_status: string; last_test_at: unknown; last_test_success: unknown }
export interface SeedRule { id: string; device_id: string; workspace_id: string; enabled: boolean; share_enabled: boolean; count: number; revision_id: string; save_state: string; runtime_apply_status?: string }
export interface SeedPoint { id: string; device_id: string; enabled: boolean; last_read_at: unknown; last_value: unknown }
export function assertSeedOrigin(baseURL: string, port: number): void {
  if (!Number.isInteger(port) || port < 1 || port > 65535 || baseURL !== `http://127.0.0.1:${port}`) throw new Error('Seeding requires the exact owned loopback origin');
}
export function assertDevice(device: SeedDevice, index: number): void {
  if (device.id !== `perf-device-${index}` || device.name !== `Synthetic offline ${index}` || device.protocol !== 'modbus_tcp' ||
    device.status !== 'draft' || device.running !== false || device.availability_status !== 'available' ||
    device.last_test_at != null || device.last_test_success != null) throw new Error('Unexpected or active seed device');
}
export function assertRule(rule: SeedRule, workspaceId: string, index: number): void {
  if (rule.id !== `perf-rule-${index}` || rule.device_id !== `perf-device-${index}` || rule.workspace_id !== workspaceId ||
    rule.enabled !== false || rule.share_enabled !== false || rule.count !== 100 || !rule.revision_id || rule.save_state !== 'saved' ||
    (rule.runtime_apply_status && !['not_running', 'aligned', 'applied'].includes(rule.runtime_apply_status))) throw new Error(`Unexpected or active seed rule ACK: ${JSON.stringify({ id: rule.id, workspace: rule.workspace_id, enabled: rule.enabled, share: rule.share_enabled, count: rule.count, save: rule.save_state, runtime: rule.runtime_apply_status })}`);
}
export function assertDisabledSources(workspaceId: string, devices: SeedDevice[], rules: SeedRule[], points: SeedPoint[]): void {
  if (!Array.isArray(devices) || devices.length !== 20 || !Array.isArray(rules) || rules.length !== 20 || !Array.isArray(points) || points.length !== 2000) throw new Error('Seed count mismatch');
  for (const list of [devices, rules, points]) if (new Set(list.map(item => item.id)).size !== list.length) throw new Error('Duplicate seed identity');
  for (let i = 1; i <= 20; i++) {
    const device = devices.find(item => item.id === `perf-device-${i}`);
    const rule = rules.find(item => item.id === `perf-rule-${i}`);
    if (!device || !rule) throw new Error('Missing seed identity');
    assertDevice(device, i); assertRule(rule, workspaceId, i);
    if (points.filter(point => point.device_id === device.id).length !== 100) throw new Error('Seed point ownership/count mismatch');
  }
  if (points.some(point => !point.id || point.enabled !== false || point.last_read_at != null || point.last_value != null)) throw new Error('Active or previously read seed point');
}
