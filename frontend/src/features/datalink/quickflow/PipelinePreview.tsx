import { useTranslation } from 'react-i18next';
import type { RuntimeStreamConnectionState } from '../../../types/datalink';
import type { QuickFlowPipelineNode } from './quickflowTypes';

interface PipelinePreviewProps {
  runtimeState: RuntimeStreamConnectionState;
  nodes: QuickFlowPipelineNode[];
  blocker: string;
}

function getStateLabel(t: ReturnType<typeof useTranslation>['t'], state: RuntimeStreamConnectionState) {
  if (state === 'connected') return t('quickflow.runtime.connected');
  if (state === 'connecting') return t('quickflow.runtime.connecting');
  if (state === 'error') return t('quickflow.runtime.error');
  return t('quickflow.runtime.disconnected');
}

export function PipelinePreview({ runtimeState, nodes, blocker }: PipelinePreviewProps) {
  const { t } = useTranslation();

  return (
    <aside
      style={{
        padding: 22,
        borderRadius: 24,
        background: 'linear-gradient(180deg, rgba(18,22,34,0.96), rgba(10,12,18,0.98))',
        border: '1px solid rgba(122,143,216,0.16)',
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', gap: 12, alignItems: 'center' }}>
        <div>
          <div style={{ color: '#9ba6c9', fontSize: 12 }}>{t('quickflow.pipeline.kicker')}</div>
          <h2 style={{ margin: '8px 0 0', fontSize: 24 }}>{t('quickflow.pipeline.title')}</h2>
        </div>
        <span style={{ padding: '6px 10px', borderRadius: 999, background: 'rgba(93,124,255,0.12)', color: '#dfe7ff', fontSize: 12 }}>
          {getStateLabel(t, runtimeState)}
        </span>
      </div>
      <div style={{ display: 'grid', gap: 18, marginTop: 18 }}>
        {nodes.map((node, index) => {
          const accent = node.status === 'ok' ? '#3dd68c' : node.status === 'warn' ? '#ffbf47' : '#ff5c7c';
          return (
            <article key={node.id} style={{ position: 'relative', padding: '16px 18px', borderRadius: 20, background: 'rgba(16, 20, 32, 0.96)', border: '1px solid rgba(122,143,216,0.14)' }}>
              {index < nodes.length - 1 ? (
                <span style={{ position: 'absolute', left: '50%', bottom: -28, width: 2, height: 28, background: 'linear-gradient(180deg, rgba(93,124,255,0.1), rgba(93,124,255,0.7))', transform: 'translateX(-50%)' }} />
              ) : null}
              <strong style={{ display: 'flex', alignItems: 'center', gap: 10, fontSize: 17 }}>
                <span style={{ width: 10, height: 10, borderRadius: '50%', background: accent, boxShadow: `0 0 14px ${accent}` }} />
                {node.title}
              </strong>
              <div style={{ marginTop: 8, color: '#9ba6c9' }}>{node.detail}</div>
              <div style={{ marginTop: 8, color: '#dbe3ff', fontSize: 13 }}>{node.metric}</div>
              <div style={{ marginTop: 12, height: 10, borderRadius: 999, background: 'rgba(255,255,255,0.06)', overflow: 'hidden' }}>
                <span
                  style={{
                    display: 'block',
                    width: node.animated ? '100%' : '34%',
                    height: '100%',
                    background: node.animated
                      ? 'linear-gradient(90deg, rgba(93,124,255,0.2), rgba(93,124,255,0.95), rgba(155,92,246,0.9))'
                      : 'rgba(255,92,124,0.3)',
                  }}
                />
              </div>
            </article>
          );
        })}
      </div>
      <div style={{ marginTop: 20, padding: 16, borderRadius: 18, background: 'rgba(255, 92, 124, 0.08)', border: '1px solid rgba(255,92,124,0.16)' }}>
        <strong style={{ display: 'block', fontSize: 16 }}>{t('quickflow.pipeline.blockerTitle')}</strong>
        <div style={{ marginTop: 8, color: '#ffcad5' }}>{blocker}</div>
      </div>
    </aside>
  );
}
