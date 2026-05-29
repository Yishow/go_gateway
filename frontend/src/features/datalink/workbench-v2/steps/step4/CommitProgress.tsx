import { useTranslation } from 'react-i18next';
import type { CommitLog } from '../../state/types';

/**
 * CommitProgress 元件屬性
 */
interface CommitProgressProps {
  logs: CommitLog[];
  status: 'idle' | 'committing' | 'success' | 'failed';
}

/**
 * 提交進度動畫日誌列表元件
 * 落地需求：「Commit sequence and animation」之動畫日誌與進度狀態
 */
export function CommitProgress({ logs, status }: CommitProgressProps) {
  const { t } = useTranslation('workbench-v2');

  const isCommitting = status === 'committing';

  return (
    <div className="bg-gray-900/10 border border-gray-800 rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm">
      <div className="space-y-4 flex-1">
        {/* 標題與動畫載入狀態 */}
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-white">
              {t('step4.progress_title', '部署與套用設定')}
            </h3>
            <p className="text-xs text-gray-500 mt-1">
              {t('step4.progress_subtitle', '正在閘道端執行 API 建置命令序列...')}
            </p>
          </div>
          {isCommitting && (
            <div className="flex space-x-1" data-testid="pulse-loader">
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
            </div>
          )}
        </div>

        {/* Log 滾動區域 */}
        <div className="border border-gray-800/60 rounded-xl bg-gray-950/40 p-4 h-[240px] overflow-y-auto space-y-2.5 font-mono text-xs scrollbar-thin scrollbar-thumb-gray-800">
          {logs.map((log, index) => (
            <div
              key={index}
              className="flex items-start justify-between gap-3 animate-[slideIn_0.2s_ease-out] text-gray-300"
            >
              <div className="flex items-start gap-2.5">
                {/* Emerald Check Icon */}
                <span className="text-emerald-500 font-bold select-none" aria-hidden="true">
                  ✓
                </span>
                <div>
                  <span className="text-gray-200 font-semibold block">{log.label}</span>
                  <span className="text-[10px] text-gray-500 block leading-tight mt-0.5">{log.detail}</span>
                </div>
              </div>
              <span className="px-1.5 py-0.5 rounded bg-emerald-500/10 text-emerald-400 font-bold text-[10px] border border-emerald-500/20 select-none">
                200
              </span>
            </div>
          ))}

          {logs.length === 0 && (
            <div className="h-full flex items-center justify-center text-gray-500 text-xs">
              {t('step4.waiting_logs', '準備發送指令...')}
            </div>
          )}
        </div>
      </div>

      {/* 底部 pulse-dot */}
      {isCommitting && (
        <div className="mt-4 flex items-center gap-2 text-xs text-blue-400 select-none bg-blue-500/5 border border-blue-500/10 rounded-lg p-2.5 animate-pulse">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
          </span>
          <span>{t('step4.executing_next_command', '正在執行下一個指令…')}</span>
        </div>
      )}
    </div>
  );
}
