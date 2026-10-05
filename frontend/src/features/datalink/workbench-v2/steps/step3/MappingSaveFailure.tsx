import { useTranslation } from 'react-i18next';
import type { Mapping } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';

interface Props {
  pointId: string;
  mapping: Mapping;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

const MAPPING_CODES = ['workspace_mapping_conflict', 'workspace_mapping_not_found', 'workspace_mapping_save_failed'];

/** Safe mapping failure copy and an explicit queue continuation. */
export function MappingSaveFailure({ pointId, mapping, dispatch }: Props) {
  const { t } = useTranslation('mapping-errors');
  const detail = mapping.save_error_detail;
  const code = MAPPING_CODES.includes(detail?.code ?? '') ? detail?.code : undefined;
  const requestId = detail?.requestId ?? mapping.save_error?.match(/Request ID:\s*([A-Za-z0-9_-]{1,128})/i)?.[1];
  const action = code === 'workspace_mapping_conflict' ? 'review_source' : code === 'workspace_mapping_not_found' ? 'reload' : 'retry';
  const removing = detail?.operation === 'delete';
  return (
    <>
      <span> · {t(removing ? 'remove_failed' : code ?? 'generic')}</span>
      <span> · {t(action)}</span>
      {code && <span> · {code}</span>}
      {detail?.status && <span> · HTTP {detail.status}</span>}
      {requestId && <span data-testid={`mapping-save-request-id-${pointId}`}> · {t('request_id')}: {requestId}</span>}
      <button type="button" className="ml-2 underline" data-testid={`mapping-save-retry-${pointId}`} onClick={(event) => {
        event.stopPropagation();
        dispatch({ type: 'retryMappingSave', pointId });
      }}>
        {t(removing ? 'retry_remove' : 'retry_button')}
      </button>
    </>
  );
}
