import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Rule, Device, ShareLayout } from '../../state/types';
import { useDeviceColor, getColorTheme, DEVICE_COLORS } from '../../state/deviceColors';
import { Icon } from '../../components';

export interface MergedPointTableProps {
  points: Point[];
  rules: Rule[];
  devices: Device[];
  conflictAddrs: Set<string>;
  shareLayouts: Record<string, ShareLayout | null>;
  invalidRules?: Rule[];
  onContinue: () => void;
}

interface MergedPointRowProps {
  point: Point;
  rules: Rule[];
  devices: Device[];
  isConflict: boolean;
  shareAddr?: number;
}

/**
 * 合併點位表的單一資料列元件
 * 
 * 落地設計決策：規避迴圈呼叫 Hook，封裝單列渲染與 useDeviceColor / getColorTheme 呼叫。
 */
const MergedPointRow: React.FC<MergedPointRowProps> = ({
  point,
  rules,
  devices,
  isConflict,
  shareAddr,
}) => {
  // 取得裝置與規則主題色
  const devTheme = useDeviceColor(point.device_id);
  const ruleIdx = rules.findIndex((r) => r.id === point.rule_id);
  const ruleTheme = getColorTheme(DEVICE_COLORS[ruleIdx === -1 ? 0 : ruleIdx % 6]);

  return (
    <tr
      className={`transition-colors hover:bg-slate-900/30 ${
        point.skipped ? 'opacity-40 bg-slate-950/10' : ''
      } ${isConflict ? 'bg-red-500/5' : ''}`}
      data-testid={`point-row-${point.id}`}
    >
      {/* 1. 設備來源 */}
      <td className="p-3">
        <div className="flex items-center gap-1.5">
          <span className={`w-1.5 h-1.5 rounded-full ${devTheme.solid}`} />
          <span className={`font-sans truncate max-w-[120px] ${devTheme.text}`}>
            {devices.find((d) => d.id === point.device_id)?.name || '未知裝置'}
          </span>
        </div>
      </td>

      {/* 2. 所屬規則 */}
      <td className="p-3">
        <span className={`font-sans truncate max-w-[120px] ${ruleTheme.text}`}>
          {point.rule_name}
        </span>
      </td>

      {/* 3. 點位名稱 */}
      <td className="p-3 text-slate-200">{point.name}</td>

      {/* 4. 暫存器位址 */}
      <td className="p-3">
        {isConflict ? (
          <span
            className="text-red-400 font-semibold flex items-center gap-1"
            data-testid={`conflict-addr-${point.id}`}
          >
            <Icon name="alert" className="w-3.5 h-3.5" />
            {point.address}
          </span>
        ) : (
          <span>{point.address}</span>
        )}
      </td>

      {/* 5. 功能碼 */}
      <td className="p-3 text-slate-400">{point.function}</td>

      {/* 6. 資料型別 */}
      <td className="p-3 text-slate-400">{point.data_type}</td>

      {/* 7. 線性轉換 */}
      <td className="p-3 text-slate-400">
        {point._rule_scale === 1 && point._rule_offset === 0 ? (
          <span className="text-slate-600">無</span>
        ) : (
          <span>
            x * {point._rule_scale}
            {point._rule_offset >= 0
              ? ` + ${point._rule_offset}`
              : ` - ${Math.abs(point._rule_offset)}`}
          </span>
        )}
      </td>

      {/* 8. 狀態標記 */}
      <td className="p-3">
        {point.skipped ? (
          <span
            className="px-1.5 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-500 text-[10px]"
            data-testid="status-skipped"
          >
            跳過
          </span>
        ) : isConflict ? (
          <span
            className="px-1.5 py-0.5 rounded bg-red-500/10 border border-red-500/20 text-red-400 text-[10px]"
            data-testid="status-conflict"
          >
            衝突
          </span>
        ) : shareAddr !== undefined ? (
          <span
            className="px-1.5 py-0.5 rounded bg-blue-500/10 border border-blue-500/20 text-blue-400 text-[10px]"
            data-testid="status-share"
          >
            Share:{shareAddr}
          </span>
        ) : (
          <span className="text-slate-600">—</span>
        )}
      </td>
    </tr>
  );
};

/**
 * 合併點位表元件
 * 
 * 落地設計決策：「Merged point table summary」與「衝突偵測」
 */
export const MergedPointTable: React.FC<MergedPointTableProps> = ({
  points,
  rules,
  devices,
  conflictAddrs,
  shareLayouts,
  invalidRules = [],
  onContinue,
}) => {
  const { t } = useTranslation('workbench-v2');
  // 1. 計算啟用點位總數
  const totalEnabled = points.filter((p) => !p.skipped && p.enabled).length;

  // 2. 計算全體點位在 Modbus Share 中的轉發地址
  const shareAddresses: Record<string, number> = {};
  const ruleEnabledCounters: Record<string, number> = {};

  points.forEach((p) => {
    if (p.skipped || !p.enabled) return;
    const layout = shareLayouts[p.rule_id];
    if (!layout) return;
    const count = ruleEnabledCounters[p.rule_id] || 0;
    shareAddresses[p.id] = layout.start + count * layout.stride;
    ruleEnabledCounters[p.rule_id] = count + 1;
  });

  return (
    <div
      className="flex flex-col gap-4 p-5 rounded-xl border border-slate-800 bg-slate-900/20"
      data-testid="merged-point-table-container"
    >
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-slate-200">
          全域合併點位總覽 (共 {points.length} 點)
        </h3>
        <span className="text-xs text-slate-500 font-mono">
          已啟用 {totalEnabled} · 已略過 {points.length - totalEnabled} · 衝突 {conflictAddrs.size}
        </span>
      </div>

      {/* Sticky Header Table Container */}
      <div className="max-h-[300px] overflow-y-auto border border-slate-800 rounded-lg bg-slate-950/40">
        <table className="w-full text-left border-collapse text-xs">
          <thead className="sticky top-0 bg-slate-900/90 backdrop-blur-sm z-10 border-b border-slate-800 text-slate-400 font-medium">
            <tr>
              <th className="p-3">設備來源</th>
              <th className="p-3">所屬規則</th>
              <th className="p-3">點位名稱</th>
              <th className="p-3">暫存器位址</th>
              <th className="p-3">功能碼</th>
              <th className="p-3">資料型別</th>
              <th className="p-3">線性轉換</th>
              <th className="p-3">狀態標記</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60 text-slate-300 font-mono">
            {points.length === 0 ? (
              <tr>
                <td colSpan={8} className="p-8 text-center text-slate-500">
                  暫無衍生點位，請配置並啟用接入規則。
                </td>
              </tr>
            ) : (
              points.map((point) => (
                <MergedPointRow
                  key={point.id}
                  point={point}
                  rules={rules}
                  devices={devices}
                  isConflict={conflictAddrs.has(point.address)}
                  shareAddr={shareAddresses[point.id]}
                />
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* 表尾與繼續按鈕 */}
      <div className="flex items-center justify-between pt-2 border-t border-slate-800/60 mt-2">
        <div className="text-xs text-slate-400">
          {conflictAddrs.size > 0 ? (
            <span className="text-red-400 flex items-center gap-1">
              <Icon name="alert" className="w-4 h-4" />
              位址存在衝突，請先排除衝突才能繼續。
            </span>
          ) : totalEnabled === 0 ? (
            <span className="text-slate-500">請啟用至少一個點位。</span>
          ) : (
            <span className="text-slate-400">
              Derivation 完畢，點位將作為 Step 3 映射配置的來源。
            </span>
          )}
        </div>

        {invalidRules.length > 0 && (
          <div
            className="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-200"
            role="alert"
            data-testid="invalid-rules-warning"
          >
            <div>{t('step2.summary.invalid_address', 'Address is invalid for the selected protocol.')}</div>
            <ul className="mt-1 list-disc pl-4">
              {invalidRules.map((rule) => (
                <li key={rule.id} data-testid={`invalid-rule-warning-${rule.id}`}>
                  {rule.name}
                </li>
              ))}
            </ul>
          </div>
        )}

        <button
          type="button"
          onClick={onContinue}
          disabled={totalEnabled === 0 || conflictAddrs.size > 0 || invalidRules.length > 0}
          className="bg-blue-600 hover:bg-blue-500 active:bg-blue-700 text-white text-xs font-semibold px-4 py-2 rounded-lg transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
          data-testid="continue-step3-btn"
        >
          繼續到映射
        </button>
      </div>
    </div>
  );
};
