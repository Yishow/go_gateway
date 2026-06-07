import { useTranslation } from 'react-i18next';
import type { DbConnector } from '../../state/types';

type DeliveryTruthStripProps = {
  connector: DbConnector;
};

type DeliveryTruthItem = {
  label: string;
  value: string;
  status: 'success' | 'failed' | 'unknown';
  at?: string | null;
  error?: string;
};

const STATUS_CLASS: Record<DeliveryTruthItem['status'], string> = {
  success: 'text-emerald-300',
  failed: 'text-rose-300',
  unknown: 'text-slate-400',
};

export function DeliveryTruthStrip({ connector }: DeliveryTruthStripProps) {
  const { t } = useTranslation('workbench-v2');

  const items: DeliveryTruthItem[] = [
    {
      label: t('step4.delivery_readiness'),
      value: t(readinessKey(connector.status)),
      status: readinessStatus(connector.status),
    },
    {
      label: t('step4.delivery_schema'),
      value: t(outcomeKey('schema', connector.last_schema_ensure_status)),
      status: outcomeStatus(connector.last_schema_ensure_status),
      at: connector.last_schema_ensure_at,
      error: connector.last_schema_ensure_error,
    },
    {
      label: t('step4.delivery_write'),
      value: t(outcomeKey('write', connector.last_write_status)),
      status: outcomeStatus(connector.last_write_status),
      at: connector.last_write_at,
      error: connector.last_write_error,
    },
    {
      label: t('step4.delivery_flush'),
      value: t(outcomeKey('flush', connector.last_flush_status)),
      status: outcomeStatus(connector.last_flush_status),
      at: connector.last_flush_at,
      error: connector.last_flush_error,
    },
  ];

  return (
    <div
      data-testid="database-delivery-truth"
      className="border-y border-gray-800/70 py-3"
    >
      <div className="mb-2 flex items-center justify-between gap-3">
        <div className="text-[11px] font-semibold uppercase tracking-wide text-slate-400">
          {t('step4.delivery_truth_title')}
        </div>
        <div className="text-[11px] text-slate-500">
          {t('step4.delivery_truth_hint')}
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-x-4 gap-y-3">
        {items.map((item) => (
          <div key={item.label} className="min-w-0">
            <div className="text-[11px] text-slate-500">{item.label}</div>
            <div className={`mt-0.5 text-sm font-medium ${STATUS_CLASS[item.status]}`}>
              {item.value}
            </div>
            {item.at && (
              <div className="mt-0.5 truncate text-[11px] text-slate-500">
                {formatDeliveryTime(item.at)}
              </div>
            )}
            {item.status === 'failed' && item.error && (
              <div className="mt-0.5 break-words text-[11px] text-rose-200">
                {item.error}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

function readinessKey(status: string): string {
  switch (status) {
    case 'ready':
      return 'step4.delivery_readiness_ready';
    case 'unreachable':
      return 'step4.delivery_readiness_unreachable';
    case 'auth_failed':
      return 'step4.delivery_readiness_auth_failed';
    case 'error':
      return 'step4.delivery_readiness_error';
    default:
      return 'step4.delivery_readiness_unknown';
  }
}

function readinessStatus(status: string): DeliveryTruthItem['status'] {
  switch (status) {
    case 'ready':
      return 'success';
    case 'unreachable':
    case 'auth_failed':
    case 'error':
      return 'failed';
    default:
      return 'unknown';
  }
}

function outcomeKey(kind: 'schema' | 'write' | 'flush', status?: string): string {
  if (status === 'success') {
    return `step4.delivery_${kind}_success`;
  }
  if (status === 'failed') {
    return `step4.delivery_${kind}_failed`;
  }
  return `step4.delivery_${kind}_unknown`;
}

function outcomeStatus(status?: string): DeliveryTruthItem['status'] {
  if (status === 'success') {
    return 'success';
  }
  if (status === 'failed') {
    return 'failed';
  }
  return 'unknown';
}

function formatDeliveryTime(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
}
