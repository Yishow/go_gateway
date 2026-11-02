import * as React from 'react';
import type { StudioLocatorEntryRead } from '../../../../types/studioV2LocatorEntry';
import type { WorkbenchV2State } from '../state/types';
/** Preserve locator consumption independently of local step presentation. Returning
 * to already-loaded content never rehydrates or clears its unsent state. */
export function useLoadedLocatorReturn({ locatorEntry, locatorSelection, locatorDraftBlocked, state, setView, setCurrent, setViewDirect, setCurrentDirect, selectRule }: {
  locatorEntry?: StudioLocatorEntryRead; locatorSelection?: Extract<StudioLocatorEntryRead['selection'], { state: 'selected' }>;
  locatorDraftBlocked?: boolean; state: WorkbenchV2State;
  setView: (view: 'flow' | 'settings') => void; setCurrent: (step: 1 | 2 | 3 | 4) => void;
  setViewDirect: (view: 'flow' | 'settings') => void; setCurrentDirect: (step: 1 | 2 | 3 | 4) => void; selectRule: (id: string) => void;
}) {
  const consumedLocator = React.useRef<string | undefined>(undefined);
  React.useLayoutEffect(() => {
    const locator = locatorEntry?.locator, key = locatorSelection && locator ? JSON.stringify([locator.workspace_id, locator.step, locator.device_id, locator.rule_id, locator.mapping_id, locator.group_id, locator.group_revision]) : undefined;
    if (!locatorSelection || !key || consumedLocator.current === key) { if (!key) consumedLocator.current = undefined; return; }
    consumedLocator.current = key;
    setView('flow'); setCurrent(locatorSelection.target.step);
    if (locatorSelection.target.ruleId) selectRule(locatorSelection.target.ruleId);
  }, [locatorEntry, locatorSelection, setView, setCurrent, selectRule]);
  // A rejected foreign selection retains the original loaded form, not a hidden
  // dirty owner. This restores presentation only; it never rehydrates or clears edits.
  React.useLayoutEffect(() => {
    if (!locatorDraftBlocked || !locatorSelection) return;
    setViewDirect('flow'); setCurrentDirect(locatorSelection.target.step);
    if (locatorSelection.target.ruleId) selectRule(locatorSelection.target.ruleId);
  }, [locatorDraftBlocked, locatorSelection, setViewDirect, setCurrentDirect, selectRule]);
  const returnToLoadedDevice = (deviceId: string) => {
    if (!locatorSelection || locatorSelection.target.deviceId !== deviceId || !state.devices.some(device => device.id === deviceId)) return false;
    setViewDirect('flow'); setCurrentDirect(locatorSelection.target.step);
    if (locatorSelection.target.ruleId) selectRule(locatorSelection.target.ruleId);
    return true;
  };
  return returnToLoadedDevice;
}
