import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { studioV2WorkspaceDatabaseAPI } from '../../../../../services/studioV2WorkspaceDatabase';
import type { DbConnector } from '../../state/types';

interface SchemaSetupSectionProps {
  connector: DbConnector;
  schemaActionsDisabled?: boolean;
  schemaPreviewSignature?: string;
}

export function SchemaSetupSection({
  connector,
  schemaActionsDisabled = false,
  schemaPreviewSignature,
}: SchemaSetupSectionProps) {
  const { t } = useTranslation('workbench-v2');
  const previewStateKey = schemaPreviewSignature ?? `${connector.kind}|${connector.schema}|${connector.table}`;
  const [schemaPhase, setSchemaPhase] = useState<'idle' | 'previewing' | 'previewed' | 'creating' | 'done' | 'error'>('idle');
  const [schemaStatements, setSchemaStatements] = useState<string[]>([]);
  const [schemaError, setSchemaError] = useState<string | null>(null);
  const [schemaExecuted, setSchemaExecuted] = useState(0);
  const schemaBusy = schemaPhase === 'previewing' || schemaPhase === 'creating';
  const schemaControlsDisabled = schemaBusy || schemaActionsDisabled;

  useEffect(() => {
    setSchemaPhase('idle');
    setSchemaStatements([]);
    setSchemaError(null);
    setSchemaExecuted(0);
  }, [previewStateKey]);

  if (connector.kind === 'sqlite') {
    return null;
  }

  const handlePreviewSchema = async () => {
    setSchemaPhase('previewing');
    setSchemaError(null);
    try {
      const result = await studioV2WorkspaceDatabaseAPI.generateSchema(true);
      setSchemaStatements(result.statements);
      setSchemaPhase('previewed');
    } catch (error) {
      setSchemaError(error instanceof Error ? error.message : String(error));
      setSchemaPhase('error');
    }
  };

  const handleCreateSchema = async () => {
    setSchemaPhase('creating');
    setSchemaError(null);
    try {
      const result = await studioV2WorkspaceDatabaseAPI.generateSchema(false);
      setSchemaStatements([]);
      setSchemaExecuted(result.executed);
      setSchemaPhase('done');
    } catch (error) {
      setSchemaError(error instanceof Error ? error.message : String(error));
      setSchemaPhase('error');
    }
  };

  return (
    <div className="space-y-2 rounded-xl border border-gray-800/60 bg-gray-950/30 p-3">
      <div className="flex items-center justify-between gap-2">
        <div>
          <div className="text-xs text-gray-300">
            {t('step4.schema_section_title', { defaultValue: '資料表結構' })}
          </div>
          <div className="mt-1 text-[11px] text-slate-500">
            {t('step4.schema_section_hint', { defaultValue: '欄位清單會依目前的 schema/table 顯示；若表不存在，可先預覽 DDL 或直接建立。' })}
          </div>
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            disabled={schemaControlsDisabled}
            onClick={handlePreviewSchema}
            className="rounded-md border border-gray-700 px-2 py-1 text-[11px] text-gray-300 transition-colors hover:bg-gray-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {schemaPhase === 'previewing'
              ? t('step4.schema_previewing', '預覽中...')
              : t('step4.schema_preview_btn', '預覽 DDL')}
          </button>
          <button
            type="button"
            disabled={schemaControlsDisabled}
            onClick={handleCreateSchema}
            className="rounded-md border border-blue-700 px-2 py-1 text-[11px] text-blue-300 transition-colors hover:bg-blue-900/40 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {schemaPhase === 'creating'
              ? t('step4.schema_creating', '建立中...')
              : t('step4.schema_create_btn', '建立資料表')}
          </button>
        </div>
      </div>

      {schemaPhase === 'previewed' && (
        schemaStatements.length === 0 ? (
          <p className="text-[11px] text-emerald-400">
            {t('step4.schema_no_changes', '資料表已存在，無需建立')}
          </p>
        ) : (
          <pre className="max-h-32 overflow-auto whitespace-pre-wrap rounded-md bg-black/40 p-2 font-mono text-[10px] text-gray-300">
            {schemaStatements.join('\n')}
          </pre>
        )
      )}

      {schemaPhase === 'done' && (
        <p className="text-[11px] text-emerald-400">
          ✓ {t('step4.schema_done', '已建立 {{count}} 項資料表結構', { count: schemaExecuted })}
        </p>
      )}

      {schemaPhase === 'error' && schemaError && (
        <p className="break-words text-[11px] text-rose-400">
          {t('step4.schema_error', '建表失敗')}: {schemaError}
        </p>
      )}
    </div>
  );
}
