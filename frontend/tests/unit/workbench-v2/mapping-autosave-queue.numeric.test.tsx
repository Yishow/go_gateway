import { render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { api, mountQueue, points, response } from './mapping-autosave-queue.harness';
import { mappingFixture } from './mapping-autosave-page.testHarness';
import { buildDefaultMapping } from '../../../src/features/datalink/workbench-v2/state/mappingDefaults';
import { MappingPreviewCells } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells';
import { mappingAPI } from '../../../src/services/datalink';
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/datalink', () => ({ mappingAPI: { preview: vi.fn() } }));
beforeEach(() => vi.resetAllMocks());
it.each([
  { source: 'uint64', scale: 1, offset: 0, target: 'uint64', raw: '9007199254740993', final: '9007199254740993', label: '9007199254740993' },
  { source: 'int16', scale: 0.5, offset: 10, target: 'float64', raw: 243, final: 131.5, label: '131.50' },
])('preserves $source defaults in the saved request and exact server preview', async ({ source, scale, offset, target, raw, final, label }) => {
  const point = { ...points[0], data_type: source, _rule_scale: scale, _rule_offset: offset };
  const record = mappingFixture({ id: 'mapping-a0', ...buildDefaultMapping(point, 0).local_value });
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  vi.mocked(mappingAPI.preview).mockResolvedValue({ raw_value: raw, final_value: final, step_results: [
    { step_index: 1, step_type: 'decode', input_value: raw, output_value: raw, error: '' },
    { step_index: 2, step_type: 'scale', input_value: raw, output_value: final, error: '' },
    { step_index: 3, step_type: 'cast', input_value: final, output_value: final, error: '' },
  ] });
  const hook = await mountQueue([point], [record]);
  hook.send({ type: 'setAllMappingsEnabled', enabled: true });
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  expect(api.update.mock.calls[0][1]).toMatchObject({ target_type: target, scale, offset });
  render(<table><tbody><tr><MappingPreviewCells point={point} mapping={hook.result.current.state.mappings.a0} rawValue={raw} /></tr></tbody></table>);
  await waitFor(() => expect(screen.getByTestId('preview-final-a0')).toHaveTextContent(label));
  expect(mappingAPI.preview).toHaveBeenCalledWith({ workspace_id: 'workspace-1', raw_value: raw, transform_pipeline: [
    { type: 'decode', order: 1, params: { data_type: source } },
    { type: 'scale', order: 2, params: { scale, offset } },
    { type: 'cast', order: 3, params: { target_type: target } },
  ] });
});
