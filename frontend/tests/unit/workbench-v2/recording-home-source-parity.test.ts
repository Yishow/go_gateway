import fixtureJSON from '../../../../testdata/recording_home/source_projection.json?raw';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { readRecordingSourceDiscovery } from '@/features/datalink/workbench-v2/recording-editor/recordingSourceDiscovery';
import type { RecordingHomeScope } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { pointAPI } from '@/services/datalink';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn() } }));

interface SourceRound {
  workspace: Awaited<ReturnType<typeof studioV2WorkspaceAPI.get>>;
  devices: Awaited<ReturnType<typeof studioV2WorkspaceDevicesAPI.list>>;
  rules: Awaited<ReturnType<typeof studioV2RulesAPI.list>>;
  mappings: Awaited<ReturnType<typeof studioV2MappingsAPI.list>>;
  points: Awaited<ReturnType<typeof pointAPI.list>>;
}
interface ParityCase {
  name: string;
  scope: RecordingHomeScope;
  before: SourceRound;
  after?: SourceRound;
  read_error?: { source: keyof SourceRound; round: number };
  expected: Awaited<ReturnType<typeof readRecordingSourceDiscovery>>;
}
// Shared with Go: literal inputs and independent expected outputs, not snapshots
// regenerated from the implementation. Only consumed transport fields are stored.
const fixtures = JSON.parse(fixtureJSON) as { schema_version: number; cases: ParityCase[] };

function arrange(fixture: ParityCase) {
  function read<K extends keyof SourceRound>(source: K) {
    let round = 0;
    return async (): Promise<SourceRound[K]> => {
      round++;
      if (fixture.read_error?.source === source && fixture.read_error.round === round) {
        throw new Error('password=secret private.host');
      }
      return structuredClone((round === 1 ? fixture.before : fixture.after ?? fixture.before)[source]);
    };
  }
  vi.mocked(studioV2WorkspaceAPI.get).mockImplementation(read('workspace'));
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockImplementation(read('devices'));
  vi.mocked(studioV2RulesAPI.list).mockImplementation(read('rules'));
  vi.mocked(studioV2MappingsAPI.list).mockImplementation(read('mappings'));
  vi.mocked(pointAPI.list).mockImplementation(read('points'));
}

beforeEach(() => vi.resetAllMocks());
describe('shared canonical source projection fixtures: legacy transport', () => {
  it('uses the documented literal fixture format', () => {
    expect(fixtures.schema_version).toBe(1);
    expect(new Set(fixtures.cases.map(item => item.name)).size).toBe(fixtures.cases.length);
    expect(fixtures.cases.length).toBeGreaterThan(0);
  });
  it.each(fixtures.cases)('$name', async fixture => {
    arrange(fixture);
    const observed = await readRecordingSourceDiscovery(fixture.scope);
    expect(observed).toEqual(fixture.expected);
    expect(JSON.stringify(observed)).not.toMatch(/password|private.host|private-value/);
    const expectedReads = fixture.read_error?.round === 1 ? 1 : 2;
    for (const read of [studioV2WorkspaceAPI.get, studioV2WorkspaceDevicesAPI.list,
      studioV2RulesAPI.list, studioV2MappingsAPI.list, pointAPI.list]) {
      expect(read).toHaveBeenCalledTimes(expectedReads);
    }
  });
});
