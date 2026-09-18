import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping, DbRowGroup, DbTarget } from '../../state/types';
import type { DbColumn } from '../../state/dbSchemas';

interface TargetMappingRowProps {
  point: Point;
  mapping: Mapping;
  target: DbTarget;
  rowGroups: DbRowGroup[];
  columns: DbColumn[];
  isConflicting: boolean;
  onUpdateTarget: (pointId: string, patch: Partial<DbTarget>) => void;
  disabled: boolean;
}

interface RowEditorState {
  draft: string;
  saved: string;
  remoteConflict: boolean;
}

const saveStateKeys = {
  saving: 'step4.saveStates.saving',
  saved: 'step4.saveStates.saved',
  'save-error': 'step4.saveStates.save-error',
  'draft-invalid': 'step4.saveStates.draft-invalid',
} as const;

function initialEditorState(columnName: string): RowEditorState {
  return {
    draft: columnName,
    saved: columnName,
    remoteConflict: false,
  };
}

/** 單一資料庫目標列，將輸入草稿與伺服器回覆分開保存。 */
export function TargetMappingRow({
  point,
  mapping,
  target,
  rowGroups,
  columns,
  isConflicting,
  onUpdateTarget,
  disabled,
}: TargetMappingRowProps) {
  const { t } = useTranslation('workbench-v2');
  const [editor, setEditor] = useState(() => initialEditorState(target.column_name));
  // 依提交順序保存尚未得到伺服器確認的值。
  const pendingCommitsRef = useRef<string[]>([]);
  const composingRef = useRef(false);

  useEffect(() => {
    const incoming = target.column_name ?? '';
    const saveState = target.save_state;
    const pendingCommits = pendingCommitsRef.current;
    const settledCommitIndex = pendingCommits.indexOf(incoming);
    const isCommitSettled = settledCommitIndex >= 0;
    const isRejected = saveState === 'save-error' || saveState === 'draft-invalid';

    if (isCommitSettled && (saveState === 'saved' || isRejected)) {
      // 這個值已有結果（儲存或被拒絕）；在它之前提交的值都已被它取代。
      pendingCommitsRef.current = pendingCommits.slice(settledCommitIndex + 1);
    }

    setEditor((current) => {
      const isDirty = current.draft !== current.saved;

      if (isCommitSettled) {
        if (saveState === 'saved') {
          if (current.draft === incoming) {
            return initialEditorState(incoming);
          }
          return { ...current, saved: incoming, remoteConflict: false };
        }
        return current;
      }

      if (!isDirty) {
        return initialEditorState(incoming);
      }

      if (incoming === current.saved) {
        return current;
      }

      if (isRejected) {
        if (incoming === current.draft) {
          return current;
        }
        return { ...current, saved: incoming, remoteConflict: true };
      }

      if (incoming === current.draft) {
        return initialEditorState(incoming);
      }

      return { ...current, saved: incoming, remoteConflict: true };
    });
  }, [target.column_name, target.save_state]);

  const commitDraft = (draftValue: string, allowConflict = false) => {
    if (disabled || !target.enabled || composingRef.current || (editor.remoteConflict && !allowConflict)) {
      return;
    }

    const nextValue = draftValue.trim();
    const currentEditor = editor;

    if (!nextValue) {
      setEditor((current) => ({ ...current, draft: current.saved, remoteConflict: false }));
      return;
    }
    if (nextValue === currentEditor.saved) {
      setEditor((current) => ({ ...current, draft: nextValue, remoteConflict: false }));
      return;
    }
    const pendingCommits = pendingCommitsRef.current;
    // 只擋與最近一次提交相同的值（Enter 後 blur 的重複觸發）；改回較早的值必須能再次送出。
    if (pendingCommits[pendingCommits.length - 1] === nextValue) {
      return;
    }

    pendingCommitsRef.current = [...pendingCommits.filter((value) => value !== nextValue), nextValue];
    onUpdateTarget(point.id, { column_name: nextValue });
  };

  const reloadRemoteValue = () => {
    setEditor((current) => initialEditorState(current.saved));
  };

  const resubmitDraft = () => {
    commitDraft(editor.draft, true);
  };

  const currentColumn = columns.find((column) => column.name === target.column_name);
  const columnType = currentColumn?.type ?? 'custom';
  const saveState = target.save_state;

  return (
    <tr
      className={`
        transition-colors duration-150
        ${isConflicting ? 'bg-red-500/5 hover:bg-red-500/10 border-l-2 border-l-red-500' : 'hover:bg-gray-800/20'}
      `}
    >
      <td className="p-3 font-mono text-xs text-gray-300">
        {mapping.tag_key}
      </td>

      <td className="p-3 text-xs text-gray-400">
        {point.name} <span className="text-gray-600 font-mono">({point.address})</span>
      </td>

      <td className="p-3 text-center text-gray-600">
        →
      </td>

      <td className="p-3">
        <div className="flex items-start gap-2">
          <div>
            <input
              type="text"
              value={editor.draft}
              list={`step4-target-column-options-${point.id}`}
              data-testid={`step4-target-column-${point.id}`}
              disabled={disabled || !target.enabled}
              onChange={(event) => {
                const nextDraft = event.target.value;
                setEditor((current) => ({ ...current, draft: nextDraft }));
              }}
              onCompositionStart={() => {
                composingRef.current = true;
              }}
              onCompositionEnd={() => {
                composingRef.current = false;
              }}
              onBlur={() => commitDraft(editor.draft)}
              onKeyDown={(event) => {
                if (event.key !== 'Enter' || composingRef.current || event.nativeEvent.isComposing) {
                  return;
                }
                event.preventDefault();
                commitDraft(editor.draft);
                event.currentTarget.blur();
              }}
              className={`
                w-44 bg-gray-950 border text-xs rounded-lg px-2 py-1.5 outline-none transition-all focus:ring-1 focus:ring-blue-500
                ${isConflicting
                  ? 'border-red-500/60 text-red-200'
                  : 'border-gray-800 text-gray-300 disabled:text-gray-500'
                }
                disabled:opacity-50 disabled:cursor-not-allowed
              `}
            />
            <datalist id={`step4-target-column-options-${point.id}`}>
              {columns.filter((column) => !column.primary_key).map((column) => (
                <option key={column.name} value={column.name} />
              ))}
            </datalist>
            {saveState && saveState !== 'idle' && (
              <div
                className={`mt-1 text-[10px] ${saveState === 'save-error' || saveState === 'draft-invalid'
                  ? 'text-rose-400'
                  : saveState === 'saving'
                    ? 'text-amber-400'
                    : 'text-emerald-400'
                  }`}
                data-testid={`step4-target-save-state-${point.id}`}
                data-save-state={saveState}
              >
                {t(saveStateKeys[saveState])}
                {saveState === 'save-error' && target.save_error && (
                  <span data-testid={`step4-target-save-error-${point.id}`}> · {target.save_error}</span>
                )}
              </div>
            )}
          </div>
          {isConflicting && (
            <span
              className="text-red-400 text-xs flex items-center gap-1 select-none"
              title={t('step4.conflict_tooltip')}
            >
              ⚠️
            </span>
          )}
        </div>
        {editor.remoteConflict && (
          <div
            className="mt-2 rounded-md border border-amber-500/30 bg-amber-500/5 p-2 text-[10px] text-amber-200"
            data-testid={`step4-target-remote-conflict-${point.id}`}
            role="alert"
          >
            <div>{t('step4.target_remote_conflict')}</div>
            <div className="mt-1 flex gap-2">
              <button
                type="button"
                disabled={disabled}
                onClick={reloadRemoteValue}
                data-testid={`step4-target-reload-remote-${point.id}`}
                className="text-amber-100 underline disabled:cursor-not-allowed disabled:opacity-50"
              >
                {t('step4.target_reload_remote')}
              </button>
              <button
                type="button"
                disabled={disabled || !target.enabled}
                onClick={resubmitDraft}
                data-testid={`step4-target-resubmit-draft-${point.id}`}
                className="text-amber-100 underline disabled:cursor-not-allowed disabled:opacity-50"
              >
                {t('step4.target_resubmit_draft')}
              </button>
            </div>
          </div>
        )}
      </td>

      <td className="p-3 text-xs text-gray-500 font-mono">
        {target.enabled ? columnType : '—'}
      </td>

      <td className="p-3">
        <select
          value={target.row_group_id ?? ''}
          disabled={disabled || !target.enabled || rowGroups.length === 0}
          data-testid={`step4-target-row-group-${point.id}`}
          onChange={(event) => onUpdateTarget(point.id, { row_group_id: event.target.value || undefined })}
          className="w-40 rounded-lg border border-gray-800 bg-gray-950 px-2 py-1.5 text-xs text-gray-300 outline-none transition-all focus:ring-1 focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <option value="">-</option>
          {rowGroups.map((group) => (
            <option key={group.id} value={group.id}>
              {group.id}
            </option>
          ))}
        </select>
      </td>

      <td className="p-3 text-right pr-6">
        <input
          type="checkbox"
          checked={target.enabled}
          disabled={disabled}
          onChange={() => onUpdateTarget(point.id, { enabled: !target.enabled })}
          className="h-3.5 w-3.5 rounded border-slate-800 bg-slate-950 text-blue-600 focus:ring-blue-500 focus:ring-offset-slate-950 cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
        />
      </td>
    </tr>
  );
}
