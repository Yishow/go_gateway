import type { ConnectRequest } from '../../types/api';

export interface QuickDraft {
  protocol: string;
  host?: string;
  port?: number;
  unitID?: number;
  station?: number;
  route?: string;
  auth?: boolean;
}

export interface ExpertDraft {
  protocol: string;
  configJson: string;
}

export type DraftType =
  | { mode: 'quick'; draft: QuickDraft }
  | { mode: 'expert'; draft: ExpertDraft };

export const gatewayAdapter = {
  quickToPayload(draft: QuickDraft): ConnectRequest {
    const protocol = draft.protocol || 'modbus-tcp';

    let defaultPort = 502;
    if (protocol.toLowerCase().includes('fatek')) {
      defaultPort = 500;
    } else if (protocol.toLowerCase().includes('mc')) {
      defaultPort = 5000;
    }

    const payload: ConnectRequest = {
      protocol,
      config: {
        host: draft.host && draft.host.trim() !== '' ? draft.host.trim() : '127.0.0.1',
        port: typeof draft.port === 'number' && !isNaN(draft.port) ? draft.port : defaultPort,
        ...(draft.route ? { route: draft.route } : {}),
        ...(draft.auth !== undefined ? { auth: draft.auth } : {}),
      },
    };

    if (protocol.toLowerCase().includes('modbus')) {
      payload.config.unitID = typeof draft.unitID === 'number' ? draft.unitID : 1;
    }

    if (protocol.toLowerCase().includes('fatek')) {
      payload.config.station = typeof draft.station === 'number' ? draft.station : 1;
    }

    return payload;
  },

  expertToPayload(draft: ExpertDraft): ConnectRequest {
    let parsedConfig: Record<string, unknown> = {};
    if (draft.configJson && draft.configJson.trim() !== '') {
      try {
        parsedConfig = JSON.parse(draft.configJson);
        if (typeof parsedConfig !== 'object' || Array.isArray(parsedConfig) || parsedConfig === null) {
           throw new Error('Config must be a JSON object');
        }
      } catch (_error) {
        throw new Error('Invalid JSON config format');
      }
    }

    return {
      protocol: draft.protocol || 'unknown',
      config: parsedConfig,
    };
  },

  payloadToDraft(payload: ConnectRequest): DraftType {
    const { protocol, config } = payload;

    if (!config || typeof config !== 'object' || Array.isArray(config) || config === null) {
      return {
        mode: 'expert',
        draft: {
          protocol: protocol || 'unknown',
          configJson: JSON.stringify(config ?? {}, null, 2),
        },
      };
    }

    // Quick fields allowed
    const quickFields = ['host', 'port', 'unitID', 'station', 'route', 'auth'];
    const configKeys = Object.keys(config);

    // If config has keys that are not part of quick form, or has nested objects, fallback to expert
    const isComplex = configKeys.some((k) => {
      if (!quickFields.includes(k)) return true;
      const val = (config as Record<string, unknown>)[k];
      return typeof val === 'object' || Array.isArray(val);
    });

    if (isComplex) {
      return {
        mode: 'expert',
        draft: {
          protocol: protocol || 'unknown',
          configJson: JSON.stringify(config, null, 2),
        },
      };
    }

    const configRecord = config as Record<string, unknown>;

    return {
      mode: 'quick',
      draft: {
        protocol: protocol || 'unknown',
        host: typeof configRecord.host === 'string' ? configRecord.host : undefined,
        port: typeof configRecord.port === 'number' ? configRecord.port : undefined,
        unitID: typeof configRecord.unitID === 'number' ? configRecord.unitID : undefined,
        station: typeof configRecord.station === 'number' ? configRecord.station : undefined,
        route: typeof configRecord.route === 'string' ? configRecord.route : undefined,
        auth: typeof configRecord.auth === 'boolean' ? configRecord.auth : undefined,
      },
    };
  },
};
