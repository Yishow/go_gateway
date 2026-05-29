import { describe, it, expect } from 'vitest';
import { mappingReducer } from '../../../src/features/datalink/workbench-v2/state/mappingReducer';
import type { WorkbenchV2State, Point, Mapping } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

const mockPoints: Point[] = [
  {
    id: 'p-01',
    device_id: 'dev-01',
    rule_id: 'rule-01',
    rule_name: 'Holding Registers',
    name: 'SENSOR_1',
    address: '40001',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 0.1,
    _rule_offset: 0,
  },
  {
    id: 'p-02',
    device_id: 'dev-01',
    rule_id: 'rule-01',
    rule_name: 'Holding Registers',
    name: 'SENSOR_2',
    address: '40002',
    data_type: 'int16',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 0.1,
    _rule_offset: 0,
  },
];

describe('mappingReducer', () => {
  describe('initMappingsForPoints', () => {
    it('should initialize mapping for new points', () => {
      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        mappings: {},
      };

      const nextState = mappingReducer(state, {
        type: 'initMappingsForPoints',
        points: mockPoints,
      });

      expect(nextState.mappings['p-01']).toBeDefined();
      expect(nextState.mappings['p-01'].tag_key).toBe('line01.temp.inlet');
      expect(nextState.mappings['p-01'].scale).toBe(0.1);
      expect(nextState.mappings['p-02']).toBeDefined();
      expect(nextState.mappings['p-02'].tag_key).toBe('line01.temp.outlet');
    });

    it('should preserve existing user edits and remove orphans', () => {
      const existingMapping: Mapping = {
        point_id: 'p-01',
        tag_key: 'custom.tag.key',
        display_name: '使用者編輯名稱',
        unit: '°F',
        target_type: 'int32',
        scale: 1.0,
        offset: 32,
        enabled: true,
      };

      const orphanMapping: Mapping = {
        point_id: 'p-orphan',
        tag_key: 'orphan.tag',
        display_name: '孤兒',
        unit: 'V',
        target_type: 'float64',
        scale: 1.0,
        offset: 0,
        enabled: true,
      };

      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        mappings: {
          'p-01': existingMapping,
          'p-orphan': orphanMapping,
        },
      };

      const nextState = mappingReducer(state, {
        type: 'initMappingsForPoints',
        points: [mockPoints[0]],
      });

      expect(nextState.mappings['p-01']).toEqual(existingMapping);
      expect(nextState.mappings['p-orphan']).toBeUndefined();
      expect(nextState.mappings['p-02']).toBeUndefined();
    });
  });

  describe('updateMapping', () => {
    it('should immutably patch mapping properties', () => {
      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        mappings: {
          'p-01': {
            point_id: 'p-01',
            tag_key: 'orig.tag',
            display_name: '原名稱',
            unit: 'C',
            target_type: 'float64',
            scale: 1.0,
            offset: 0,
            enabled: true,
          },
        },
      };

      const nextState = mappingReducer(state, {
        type: 'updateMapping',
        pointId: 'p-01',
        patch: {
          tag_key: 'new.tag',
          scale: 2.0,
        },
      });

      expect(nextState.mappings['p-01'].tag_key).toBe('new.tag');
      expect(nextState.mappings['p-01'].scale).toBe(2.0);
      expect(nextState.mappings['p-01'].display_name).toBe('原名稱');
    });

    it('should ignore patch if point mapping does not exist', () => {
      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        mappings: {},
      };

      const nextState = mappingReducer(state, {
        type: 'updateMapping',
        pointId: 'p-nonexistent',
        patch: { tag_key: 'ignore' },
      });

      expect(nextState.mappings['p-nonexistent']).toBeUndefined();
    });
  });

  describe('toggleMappingEnabled', () => {
    it('should toggle enabled state of a mapping', () => {
      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        mappings: {
          'p-01': {
            point_id: 'p-01',
            tag_key: 'tag',
            display_name: 'name',
            unit: 'U',
            target_type: 'float64',
            scale: 1,
            offset: 0,
            enabled: true,
          },
        },
      };

      const state2 = mappingReducer(state, {
        type: 'toggleMappingEnabled',
        pointId: 'p-01',
      });
      expect(state2.mappings['p-01'].enabled).toBe(false);

      const state3 = mappingReducer(state2, {
        type: 'toggleMappingEnabled',
        pointId: 'p-01',
      });
      expect(state3.mappings['p-01'].enabled).toBe(true);
    });

    it('should return original state if point mapping does not exist', () => {
      const state = { ...INITIAL_STATE, mappings: {} };
      const nextState = mappingReducer(state, {
        type: 'toggleMappingEnabled',
        pointId: 'p-nonexistent',
      });
      expect(nextState).toBe(state);
    });
  });

  describe('bulkApplyTransform', () => {
    it('should copy specified fields from source mapping to all mappings', () => {
      const state: WorkbenchV2State = {
        ...INITIAL_STATE,
        mappings: {
          'p-01': {
            point_id: 'p-01',
            tag_key: 'tag.01',
            display_name: 'n1',
            unit: 'U',
            target_type: 'int32',
            scale: 5.0,
            offset: 10,
            enabled: true,
          },
          'p-02': {
            point_id: 'p-02',
            tag_key: 'tag.02',
            display_name: 'n2',
            unit: 'U',
            target_type: 'float64',
            scale: 1.0,
            offset: 0,
            enabled: true,
          },
        },
      };

      const nextState = mappingReducer(state, {
        type: 'bulkApplyTransform',
        fromPointId: 'p-01',
        fields: ['scale', 'offset', 'target_type'],
      });

      expect(nextState.mappings['p-02'].scale).toBe(5.0);
      expect(nextState.mappings['p-02'].offset).toBe(10);
      expect(nextState.mappings['p-02'].target_type).toBe('int32');
      expect(nextState.mappings['p-02'].tag_key).toBe('tag.02');
    });

    it('should return original state if source point mapping does not exist', () => {
      const state = { ...INITIAL_STATE, mappings: {} };
      const nextState = mappingReducer(state, {
        type: 'bulkApplyTransform',
        fromPointId: 'p-nonexistent',
        fields: ['scale'],
      });
      expect(nextState).toBe(state);
    });
  });
});
