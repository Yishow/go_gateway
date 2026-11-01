import { useLayoutEffect, useRef, type RefObject } from 'react';
import { readStudioLegacyPresentation } from '../../../../services/studioRecordingView';
import type { useWorkbenchV2State } from '../state/useWorkbenchV2State';
import { useRecordingNavigation } from '../recording-editor/RecordingDraftNavigation';
/** Bounded old links select existing surfaces; they never manufacture a locator or mutate configuration. */
export function useStudioLegacyPresentation(search: string | undefined, actions: ReturnType<typeof useWorkbenchV2State>,
  root: RefObject<HTMLDivElement | null>, enabled: boolean) {
  const { request } = useRecordingNavigation();
  const { view, current } = actions.state;
  const consumed = useRef<string | undefined>(undefined);
  const focused = useRef<string | undefined>(undefined);
  const intent = enabled ? readStudioLegacyPresentation(search) : null;
  useLayoutEffect(() => {
    if (!intent || consumed.current === search) return;
    consumed.current = search;
    request(() => { actions.setView(intent.view); actions.setCurrent(intent.step); });
  }, [intent, search, actions, request]);
  useLayoutEffect(() => {
    if (!intent?.target || focused.current === search || view !== intent.view || current !== intent.step) return;
    const target = root.current?.querySelector<HTMLElement>(intent.target === 'modbus'
      ? '[data-testid^="step4-share-output-"]' : '[data-testid="database-connection-card"]');
    if (target) { focused.current = search; target.tabIndex = -1; target.focus({ preventScroll: true }); }
  }, [intent, search, view, current, root]);
}
