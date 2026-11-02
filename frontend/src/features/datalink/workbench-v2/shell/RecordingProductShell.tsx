import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { useOptionalTheme } from '@/contexts/ThemeContext';
import './recordingProductShell.css';
/** Shared presentation only: scope comes from the current owner; theme has one existing owner. */
export function RecordingProductShell({ workspaceId, tools }: { workspaceId?: string; tools?: ReactNode }) {
  const { t } = useTranslation('recording-home');
  const theme = useOptionalTheme();
  return <header className="gw-product-header" data-testid="recording-product-header">
    <div className="gw-product-brand"><span aria-hidden="true" className="gw-product-mark">◇</span><div><strong>go_gateway</strong><p>{t('product.workspace')}</p></div></div>
    <div className="gw-product-scope"><span>{t('product.scope')}</span><code title={workspaceId}>{workspaceId || t('unknown')}</code></div>
    <div className="gw-product-tools">
      {tools && <details><summary>{t('product.tools')}</summary><nav aria-label={t('navigation')}>{tools}</nav></details>}
      {theme && <label>{t('product.appearance')}<select aria-label={t('product.appearance')} value={theme.preference} onChange={event => {
        const value = event.target.value; if (value === 'system' || value === 'light' || value === 'dark') theme.setPreference(value);
      }}><option value="system">{t('product.system')}</option><option value="light">{t('product.light')}</option><option value="dark">{t('product.dark')}</option></select></label>}
      {theme && theme.saveState !== 'saved' && theme.saveState !== 'idle' && <span role={theme.saveState === 'error' ? 'alert' : 'status'}>{t(`product.theme.${theme.saveState}`)}{theme.saveState === 'error' && <button type="button" onClick={() => { void theme.persistCurrent().catch(() => undefined); }}>{t('retry')}</button>}</span>}
    </div>
  </header>;
}
