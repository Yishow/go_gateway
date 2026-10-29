import type { ReactNode } from 'react';
import type { EntityLocator } from '../../../types/runtimeRuleEvidence';
import { buildStudioEntityLocatorTarget } from '../../../services/studioEntityLocator';
import { handleStudioTargetClick } from './RuntimeSetupHelpers';

/** A validated identity-only link, not a fetch, edit, activation or mutation capability. */
export function RuntimeSettingsLink({ locator, navigateTo, children }: {
  locator: EntityLocator | null; navigateTo?: (target: string) => void; children: ReactNode;
}) {
  const target = locator && buildStudioEntityLocatorTarget(locator);
  if (!target) return null;
  const click = handleStudioTargetClick(navigateTo, target);
  return <a href={target} onClick={click ? (event) => {
    if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
    click(event);
  } : undefined} className="inline-flex rounded border border-cyan-700 px-2 py-1 text-xs text-cyan-200 focus-visible:outline focus-visible:outline-sky-400">{children}</a>;
}
