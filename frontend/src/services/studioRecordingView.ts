import { DASHBOARD_MODAL_INTENTS, DASHBOARD_SECTION_INTENTS, LEGACY_DECOMMISSION_ROUTES, WORKBENCH_COMPAT_STEPS, WORKBENCH_COMPAT_TARGETS } from '../features/datalink/legacyRoutes';
import { buildStudioEntityLocatorTarget, parseStudioEntityLocator, type StudioEntityLocatorIntent } from './studioEntityLocator';
export type StudioPresentation = 'recordings' | 'recording' | 'advanced';
export interface StudioPresentationIntent { presentation: StudioPresentation; locatorIntent: StudioEntityLocatorIntent; recordingIntent?: 'copy' }
/** Validate raw query before stripping only one recognized presentation hint. */
export function parseStudioPresentation(search: string): StudioPresentationIntent {
  const invalid = { presentation: 'advanced', locatorIntent: { kind: 'invalid' } } as const;
  if (typeof search !== 'string' || search.length > 2048 || /[;#]/.test(search)) return invalid;
  const raw = search.startsWith('?') ? search.slice(1) : search;
  try { for (const pair of raw.split('&')) for (const part of pair.split('=')) decodeURIComponent(part.replace(/\+/g, ' ')); }
  catch { return invalid; }
  const params = new URLSearchParams(raw);
  if ([...params.keys()].some(key => params.getAll(key).length !== 1)) return invalid;
  const view = params.get('view');
  if (view === 'recordings' || view === 'recording' || view === 'advanced') params.delete('view');
  // A new local draft has no persisted ID. No locator, extra command or return URL
  // may be smuggled into this bounded presentation intent.
  if (params.has('intent')) {
    if (view === 'recording' && params.get('intent') === 'new' && params.size === 1)
      return { presentation: 'recording', locatorIntent: { kind: 'absent' } };
    if (view !== 'recording' || params.get('intent') !== 'copy') return invalid;
    params.delete('intent');
    const locatorIntent = parseStudioEntityLocator(params.toString());
    return locatorIntent.kind === 'locator' && locatorIntent.locator.step === 4
      ? { presentation: 'recording', recordingIntent: 'copy', locatorIntent } : invalid;
  }
  if (legacyPresentation(params)) return { presentation: 'advanced', locatorIntent: { kind: 'absent' } };
  const locatorIntent = parseStudioEntityLocator(params.toString());
  if (view === 'recording') {
    if (locatorIntent.kind === 'locator') return { presentation: locatorIntent.locator.step === 4 ? 'recording' : 'advanced', locatorIntent };
    if (locatorIntent.kind === 'absent' && params.get('focus') === 'readiness' && /^[1-4]$/.test(params.get('step') ?? ''))
      return { presentation: 'advanced', locatorIntent };
    return invalid;
  }
  return { locatorIntent, presentation: locatorIntent.kind === 'absent' && (view === 'recordings' || view === null && params.size === 0) ? 'recordings' : 'advanced' };
}
/** Same-route navigation only, preserving all unrelated raw query bytes and hash. */
export function buildStudioPresentationTarget(search: string, hash: string, view: Exclude<StudioPresentation, 'recording'>): string | null {
  if (parseStudioPresentation(search).locatorIntent.kind === 'invalid' || (hash && !hash.startsWith('#'))) return null;
  const parts = (search.startsWith('?') ? search.slice(1) : search).split('&').filter(Boolean)
    .filter(pair => decodeURIComponent(pair.split('=')[0].replace(/\+/g, ' ')) !== 'view')
    .filter(pair => decodeURIComponent(pair.split('=')[0].replace(/\+/g, ' ')) !== 'intent');
  parts.push(`view=${view}`);
  return `/studio/v2?${parts.join('&')}${hash}`;
}

/** Navigation intent only; the destination re-reads canonical scope before editing. */
export function buildSimpleRecordingTarget(intent: { kind: 'new' } | { kind: 'edit' | 'copy'; workspaceId: string; groupId: string; groupRevision: string }, hash = ''): string | null {
  if (hash && !hash.startsWith('#')) return null;
  if (intent.kind === 'new') return `/studio/v2?view=recording&intent=new${hash}`;
  const target = buildStudioEntityLocatorTarget({ step: 4, workspace_id: intent.workspaceId, device_id: null,
    rule_id: null, mapping_id: null, group_id: intent.groupId, group_revision: intent.groupRevision });
  return target ? `${target}&view=recording${intent.kind === 'copy' ? '&intent=copy' : ''}${hash}` : null;
}

export interface StudioLegacyPresentation { view: 'flow' | 'settings'; step: 1 | 2 | 3 | 4; target?: 'modbus' | 'database' }
/** Only the already-supported legacy vocabulary is non-locator navigation. Any
 * canonical locator key, unknown key, duplicate or malformed byte stays strict. */
function legacyPresentation(params: URLSearchParams): StudioLegacyPresentation | null {
  const allowed: Record<string, readonly string[]> = { step: WORKBENCH_COMPAT_STEPS, target: WORKBENCH_COMPAT_TARGETS,
    section: DASHBOARD_SECTION_INTENTS, modal: DASHBOARD_MODAL_INTENTS, legacy: LEGACY_DECOMMISSION_ROUTES };
  if (!params.size || [...params].some(([key, value]) => !Object.hasOwn(allowed, key) || !allowed[key].includes(value) || params.getAll(key).length !== 1)) return null;
  if (params.has('legacy') && params.get('legacy') !== params.get('modal')) return null;
  const section = params.get('section'), modal = params.get('modal'), target = params.get('target');
  const explicit = params.get('step');
  const step = explicit ? WORKBENCH_COMPAT_STEPS.indexOf(explicit as typeof WORKBENCH_COMPAT_STEPS[number]) + 1
    : target ? 4 : modal === 'mappings' || modal === 'tags' ? 3 : modal === 'points' || modal === 'polling-groups' ? 2 : 1;
  return { view: section === 'settings' || modal === 'settings' ? 'settings' : 'flow', step: step as 1 | 2 | 3 | 4,
    ...(target === 'modbus' || target === 'database' ? { target } : {}) };
}
export function readStudioLegacyPresentation(search = ''): StudioLegacyPresentation | null {
  if (parseStudioPresentation(search).locatorIntent.kind !== 'absent') return null;
  const params = new URLSearchParams(search); params.delete('view');
  return legacyPresentation(params);
}
