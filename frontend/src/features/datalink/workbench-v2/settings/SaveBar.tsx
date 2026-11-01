import { useLayoutEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';

/**
 * SaveBar 元件屬性
 */
interface SaveBarProps {
  onReset: () => void;
  onSave: () => void;
}

/**
 * 系統設定底部 Sticky 儲存控制列元件
 * 落地需求：「Save bar actions」與「Sticky 儲存列定位」
 */
export function SaveBar({ onReset, onSave }: SaveBarProps) {
  const { t } = useTranslation('workbench-v2');
  const footerRef = useRef<HTMLDivElement>(null);

  useLayoutEffect(() => {
    const footer = footerRef.current;
    if (!footer) return;
    const documentStyle = footer.ownerDocument.documentElement.style;
    const previous = documentStyle.getPropertyValue('scroll-padding-bottom');
    const priority = documentStyle.getPropertyPriority('scroll-padding-bottom');
    const nonNegative = (value: number) => Number.isFinite(value) ? Math.max(0, value) : 0;
    const update = () => {
      const height = nonNegative(footer.getBoundingClientRect().height);
      const bottom = nonNegative(Number.parseFloat(window.getComputedStyle(footer).bottom));
      documentStyle.setProperty('scroll-padding-bottom', `${height + bottom + 8}px`);
    };
    update();
    const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(update);
    observer?.observe(footer);
    window.addEventListener('resize', update);
    return () => {
      observer?.disconnect();
      window.removeEventListener('resize', update);
      if (previous) documentStyle.setProperty('scroll-padding-bottom', previous, priority);
      else documentStyle.removeProperty('scroll-padding-bottom');
    };
  }, []);

  const handleResetClick = () => {
    // 彈出原生 confirm 確認
    const ok = window.confirm(
      t('settings.reset_confirm_msg', '確定要將所有系統設定重設為預設值嗎？這將會清除您目前所有的連接器設定與變更。')
    );
    if (ok) {
      onReset();
    }
  };

  return (
    <div ref={footerRef} className="sticky bottom-4 z-10 flex flex-wrap items-center justify-between gap-4 p-4 rounded-2xl border border-[var(--gw-border)] bg-[var(--gw-surface)]">
      {/* 提示說明文案 */}
      <div className="flex items-center gap-2">
        <span className="text-base select-none">💡</span>
        <span className="text-xs text-[var(--gw-muted)]">
          {t('settings.save_bar_info', '設定變更會立即套用至本地運行，並於後端服務下次重啟後持續生效。')}
        </span>
      </div>

      {/* 按鈕操作區 */}
      <div className="flex items-center gap-3">
        {/* 重設為預設 */}
        <button
          type="button"
          onClick={handleResetClick}
          className="px-4 py-2 rounded-xl text-xs font-semibold border border-transparent hover:border-[var(--gw-control)] text-[var(--gw-muted)] hover:text-[var(--gw-text)] bg-transparent hover:bg-[var(--gw-surface)] cursor-pointer transition-all"
        >
          {t('settings.reset_to_defaults_btn', '重設為預設')}
        </button>

        {/* 儲存所有設定 */}
        <button
          type="button"
          onClick={onSave}
          className="px-5 py-2.5 rounded-xl text-xs font-semibold bg-[var(--gw-primary)] hover:bg-[var(--gw-primary-hover)] text-[var(--gw-on-primary)] cursor-pointer transition-all"
        >
          {t('settings.save_all_btn', '儲存所有設定')}
        </button>
      </div>
    </div>
  );
}
