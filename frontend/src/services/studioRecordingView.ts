import { parseStudioEntityLocator, type StudioEntityLocatorIntent } from './studioEntityLocator';
export type StudioPresentation = 'recordings' | 'advanced';
export interface StudioPresentationIntent { presentation: StudioPresentation; locatorIntent: StudioEntityLocatorIntent }
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
  if (view === 'recordings' || view === 'advanced') params.delete('view');
  const locatorIntent = parseStudioEntityLocator(params.toString());
  return { locatorIntent, presentation: locatorIntent.kind === 'absent' && view === 'recordings' ? 'recordings' : 'advanced' };
}
/** Same-route navigation only, preserving all unrelated raw query bytes and hash. */
export function buildStudioPresentationTarget(search: string, hash: string, view: StudioPresentation): string | null {
  if (parseStudioPresentation(search).locatorIntent.kind === 'invalid' || (hash && !hash.startsWith('#'))) return null;
  const parts = (search.startsWith('?') ? search.slice(1) : search).split('&').filter(Boolean)
    .filter(pair => decodeURIComponent(pair.split('=')[0].replace(/\+/g, ' ')) !== 'view');
  parts.push(`view=${view}`);
  return `/studio/v2?${parts.join('&')}${hash}`;
}
