import { describe, it, expect } from 'vitest';
import { gatewayAdapter, type QuickDraft, type ExpertDraft } from '../gatewayAdapter';
import type { ConnectRequest } from '../../../types/api';

describe('gatewayAdapter', () => {
  describe('quickToPayload', () => {
    it('should fill default host and port for modbus-tcp', () => {
      const draft: QuickDraft = { protocol: 'modbus-tcp' };
      const payload = gatewayAdapter.quickToPayload(draft);
      
      expect(payload).toEqual({
        protocol: 'modbus-tcp',
        config: {
          host: '127.0.0.1',
          port: 502,
          unitID: 1
        }
      });
    });

    it('should fill default port for fatek-tcp', () => {
      const draft: QuickDraft = { protocol: 'fatek-tcp' };
      const payload = gatewayAdapter.quickToPayload(draft);
      
      expect(payload).toEqual({
        protocol: 'fatek-tcp',
        config: {
          host: '127.0.0.1',
          port: 500,
          station: 1
        }
      });
    });

    it('should fill default port for mc-protocol', () => {
      const draft: QuickDraft = { protocol: 'mc-tcp' };
      const payload = gatewayAdapter.quickToPayload(draft);
      
      expect(payload).toEqual({
        protocol: 'mc-tcp',
        config: {
          host: '127.0.0.1',
          port: 5000,
        }
      });
    });

    it('should respect user provided host and port', () => {
      const draft: QuickDraft = { protocol: 'modbus-tcp', host: '192.168.1.10', port: 1502, unitID: 5, route: '/api/v1', auth: true };
      const payload = gatewayAdapter.quickToPayload(draft);
      
      expect(payload).toEqual({
        protocol: 'modbus-tcp',
        config: {
          host: '192.168.1.10',
          port: 1502,
          unitID: 5,
          route: '/api/v1',
          auth: true
        }
      });
    });
  });

  describe('expertToPayload', () => {
    it('should parse valid JSON without changing semantics', () => {
      const draft: ExpertDraft = {
        protocol: 'modbus-rtu',
        configJson: '{\n  "port": "COM1",\n  "baudRate": 9600\n}'
      };
      
      const payload = gatewayAdapter.expertToPayload(draft);
      expect(payload).toEqual({
        protocol: 'modbus-rtu',
        config: {
          port: 'COM1',
          baudRate: 9600
        }
      });
    });

    it('should throw error on invalid JSON', () => {
      const draft: ExpertDraft = {
        protocol: 'modbus-rtu',
        configJson: '{\n  "port": "COM1",\n  "baudRate": 9600' // missing closing brace
      };
      
      expect(() => gatewayAdapter.expertToPayload(draft)).toThrow('Invalid JSON config format');
    });

    it('should handle empty string as empty object', () => {
      const draft: ExpertDraft = { protocol: 'modbus-tcp', configJson: '   ' };
      const payload = gatewayAdapter.expertToPayload(draft);
      
      expect(payload.config).toEqual({});
    });

    it('should throw error if JSON is not an object', () => {
      const draft: ExpertDraft = { protocol: 'modbus-tcp', configJson: '["a", "b"]' };
      expect(() => gatewayAdapter.expertToPayload(draft)).toThrow('Invalid JSON config format');
    });
  });

  describe('payloadToDraft', () => {
    it('should map simple payload to quick draft', () => {
      const payload: ConnectRequest = {
        protocol: 'modbus-tcp',
        config: {
          host: '192.168.0.5',
          port: 502,
          unitID: 2,
          route: '/test',
          auth: false
        }
      };
      
      const draft = gatewayAdapter.payloadToDraft(payload);
      expect(draft).toEqual({
        mode: 'quick',
        draft: {
          protocol: 'modbus-tcp',
          host: '192.168.0.5',
          port: 502,
          unitID: 2,
          station: undefined,
          route: '/test',
          auth: false
        }
      });
    });

    it('should map complex payload to expert draft clearly', () => {
      const payload: ConnectRequest = {
        protocol: 'modbus-rtu',
        config: {
          port: '/dev/ttyUSB0', // not in quick fields
          baudRate: 115200,
          unitID: 1
        }
      };
      
      const draft = gatewayAdapter.payloadToDraft(payload);
      expect(draft.mode).toBe('expert');
      if (draft.mode === 'expert') {
        expect(draft.draft.protocol).toBe('modbus-rtu');
        expect(JSON.parse(draft.draft.configJson)).toEqual({
          port: '/dev/ttyUSB0',
          baudRate: 115200,
          unitID: 1
        });
      }
    });

    it('should handle abnormal payload fallback (null config) cleanly', () => {
      const payload = {
        protocol: 'unknown-proto',
        config: null
      } as unknown as ConnectRequest;
      
      const draft = gatewayAdapter.payloadToDraft(payload);
      expect(draft.mode).toBe('expert');
      if (draft.mode === 'expert') {
        expect(draft.draft.configJson).toBe('{}');
      }
    });

    it('should map payload with nested objects to expert draft', () => {
      const payload: ConnectRequest = {
        protocol: 'mqtt',
        config: {
          host: 'broker.hivemq.com',
          port: 1883,
          tls: {
            enabled: true
          }
        }
      };

      const draft = gatewayAdapter.payloadToDraft(payload);
      expect(draft.mode).toBe('expert');
      if (draft.mode === 'expert') {
        expect(JSON.parse(draft.draft.configJson)).toEqual(payload.config);
      }
    });
  });
});
