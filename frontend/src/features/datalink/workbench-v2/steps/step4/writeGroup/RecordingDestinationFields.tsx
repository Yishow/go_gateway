import { useQuery } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import { studioLocatorReadKey } from '../../../../../../services/studioEntityLocator';
import { readRecordingDestinationCatalog } from '../../../../../../services/studioV2DatabaseCatalog';
import type { StudioV2LocatorReadScope } from '../../../../../../types/studioV2LocatorEntry';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import type { DbConnector } from '../../../state/types';
import type { ConfiguredDestinationObservation } from '../../../recording-editor/ConfiguredDestinationChooser';
import { destinationSaved } from './GroupEditorCapabilities';
import { GroupIdentityFields, type GroupIdentityFieldsProps } from './GroupIdentityFields';
import type { RecordingCatalogView } from './RecordingCatalogFields';
import { ownedWideCatalogTables } from './recordingCatalogOwnership';

/** Bind available destination choices to the currently displayed persisted identity. */
export function RecordingDestinationFields({ workspaceId, workspaceRevision, connector, group, locatorScope, configuredDestinations, ...fields }: GroupIdentityFieldsProps & {
  workspaceId: string; workspaceRevision: string; connector: DbConnector; group?: WriteGroup; locatorScope?: StudioV2LocatorReadScope; configuredDestinations?: ConfiguredDestinationObservation;
}) {
  const marker = useRef<HTMLSpanElement>(null);
  const [destinationVisible, setDestinationVisible] = useState(false);
  useEffect(() => {
    const section = marker.current?.closest('#recording-editor-destination');
    if (!section || typeof IntersectionObserver === 'undefined') {
      setDestinationVisible(true);
      return;
    }
    const observer = new IntersectionObserver(entries => setDestinationVisible(entries.some(entry => entry.isIntersecting)));
    observer.observe(section);
    return () => observer.disconnect();
  }, []);
  const scope = { workspaceId, workspaceRevision, connectorId: connector.connector_id ?? '', connectorRevision: connector.identity_revision ?? '', kind: connector.kind,
    groupId: group?.id, groupRevision: group?.revision };
  const eligible = destinationSaved(connector) && Boolean(workspaceId) && !fields.locked && fields.draft.connector_id === scope.connectorId &&
    fields.draft.connector_revision === scope.connectorRevision && (!group || Boolean(group.id && group.revision));
  const query = useQuery({ queryKey: ['studio-v2', 'database-catalog', workspaceId, workspaceRevision, scope.connectorId, scope.connectorRevision, scope.kind,
    scope.groupId ?? '', scope.groupRevision ?? '', ...studioLocatorReadKey(locatorScope)], queryFn: () => readRecordingDestinationCatalog(scope, locatorScope),
    enabled: eligible && destinationVisible, retry: false, refetchOnWindowFocus: false, staleTime: 0 });
  // A query that has not started (destination not yet on screen) is not "loading"; show plain fields until it runs.
  const started = query.isFetching || query.isFetched;
  const status: RecordingCatalogView['status'] = query.isFetching ? 'loading' : query.isError || !query.data ? 'error' : 'ready';
  const catalog: RecordingCatalogView | undefined = eligible && started ? { status, defaultSchema: status === 'ready' ? query.data!.default_schema : '',
    schemas: status === 'ready' ? query.data!.schemas : [], blockedTables: status === 'ready' ? ownedWideCatalogTables(configuredDestinations,
      { workspaceId, workspaceRevision, connector, groupId: group?.id }, query.data!.default_schema) : [], retry: () => { void query.refetch(); } } : undefined;
  return <><span ref={marker} hidden /><GroupIdentityFields {...fields} section="destination" catalog={catalog} /></>;
}
