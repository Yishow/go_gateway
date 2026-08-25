import { boundedString, normalizeTypedEnvelope, parseBoundedJson } from '../utils/safeJson';

/** SSE preview event contract exposed by the persisted mapping stream. */
export interface PreviewEvent {
  type: 'connected' | 'preview' | 'heartbeat' | 'error' | 'close';
  mapping_id?: string;
  raw_value?: number | string | boolean;
  final_value?: number | string | boolean;
  steps?: StepResult[];
  quality?: number;
  timestamp: string;
  /** Typed error metadata only; raw backend error text is intentionally dropped. */
  error?: PreviewErrorEnvelope;
  code?: string;
  action?: string;
  request_id?: string;
  retryable?: boolean;
}

export interface PreviewErrorEnvelope {
  code?: string;
  action?: string;
  request_id?: string;
  retryable?: boolean;
}

export interface StepResult {
  step_index: number;
  step_type: string;
  input: number | string | boolean;
  output: number | string | boolean;
  error?: string;
}

const PREVIEW_EVENT_TYPES = new Set<PreviewEvent['type']>([
  'connected',
  'preview',
  'heartbeat',
  'error',
  'close',
]);

function previewValue(value: unknown): number | string | boolean | undefined {
  if (typeof value === 'number' && Number.isFinite(value)) return value;
  if (typeof value === 'string') return boundedString(value);
  if (typeof value === 'boolean') return value;
  return undefined;
}

/** Narrows untrusted SSE JSON to the preview event contract. */
export function parsePreviewEvent(payload: unknown): PreviewEvent | null {
  const boundedPayload = parseBoundedJson(payload);
  if (typeof boundedPayload !== 'object' || boundedPayload === null || Array.isArray(boundedPayload)) return null;
  const record = boundedPayload as Record<string, unknown>;
  const type = record.type;
  const timestamp = boundedString(record.timestamp);
  if (typeof type !== 'string' || !PREVIEW_EVENT_TYPES.has(type as PreviewEvent['type']) || !timestamp) {
    return null;
  }

  const event: PreviewEvent = { type: type as PreviewEvent['type'], timestamp };
  const envelope = normalizeTypedEnvelope(record);
  const mappingId = boundedString(record.mapping_id);
  if (mappingId) event.mapping_id = mappingId;
  if (envelope.code) event.code = envelope.code;
  if (envelope.action) event.action = envelope.action;
  if (envelope.requestId) event.request_id = envelope.requestId;
  if (envelope.retryable !== undefined) event.retryable = envelope.retryable;
  if (envelope.code || envelope.action || envelope.requestId || envelope.retryable !== undefined) {
    event.error = {
      ...(envelope.code ? { code: envelope.code } : {}),
      ...(envelope.action ? { action: envelope.action } : {}),
      ...(envelope.requestId ? { request_id: envelope.requestId } : {}),
      ...(envelope.retryable !== undefined ? { retryable: envelope.retryable } : {}),
    };
  }

  const rawValue = previewValue(record.raw_value);
  const finalValue = previewValue(record.final_value);
  if (rawValue !== undefined) event.raw_value = rawValue;
  if (finalValue !== undefined) event.final_value = finalValue;
  if (typeof record.quality === 'number' && Number.isFinite(record.quality)) event.quality = record.quality;
  if (Array.isArray(record.steps)) {
    event.steps = record.steps.flatMap((step): StepResult[] => {
      if (typeof step !== 'object' || step === null) return [];
      const stepRecord = step as Record<string, unknown>;
      const stepIndex = stepRecord.step_index;
      const stepType = boundedString(stepRecord.step_type);
      const input = previewValue(stepRecord.input);
      const output = previewValue(stepRecord.output);
      if (typeof stepIndex !== 'number' || !Number.isInteger(stepIndex) || stepIndex < 0 || !stepType || input === undefined || output === undefined) {
        return [];
      }
      const normalized: StepResult = { step_index: stepIndex, step_type: stepType, input, output };
      return [normalized];
    });
  }
  return event;
}
