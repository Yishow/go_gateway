import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { studioV2WorkspaceDatabaseAPI } from '../../../../../services/studioV2WorkspaceDatabase';
import { getSafeErrorMessage, type SafeErrorMessage } from '@/utils/typedErrors';
import type { DbConnector } from '../../state/types';
import {
  beginSchemaOperation,
  isExplicitSchemaHTTPFailure,
  isSchemaOperationScopeCurrent,
  settleSchemaOperation,
  useSchemaOperation,
} from './schemaOperationState';

interface SchemaSetupSectionProps {
  connector: DbConnector;
  workspaceId?: string;
  schemaActionsDisabled?: boolean;
  schemaPreviewSignature?: string;
}

type LocalSchemaPhase = 'idle' | 'previewing' | 'previewed' | 'creating' | 'done' | 'error' | 'unknown';

interface CredentialSnapshot {
  password?: string;
  clearPassword?: boolean;
  passwordRequired?: boolean;
}

function credentialSnapshot(connector: DbConnector): CredentialSnapshot {
  return {
    password: connector.password,
    clearPassword: connector.clear_password,
    passwordRequired: connector.password_required,
  };
}

/** Schema 預覽與建表操作；操作鎖跨過短暫 remount，避免重複送出未知結果。 */
export function SchemaSetupSection({
  connector,
  workspaceId,
  schemaActionsDisabled = false,
  schemaPreviewSignature,
}: SchemaSetupSectionProps) {
  const { t } = useTranslation('workbench-v2');
  const workspaceKey = workspaceId ?? connector.workspace_id ?? '__unscoped__';
  const [credentialGeneration, setCredentialGeneration] = useState(0);
  const previousCredentialRef = useRef<CredentialSnapshot>(credentialSnapshot(connector));
  const mountedRef = useRef(true);
  const basePreviewStateKey = schemaPreviewSignature ?? `${connector.kind}|${connector.schema}|${connector.table}`;

  const { password, clear_password: clearPassword, password_required: passwordRequired } = connector;
  useLayoutEffect(() => {
    const nextCredential: CredentialSnapshot = { password, clearPassword, passwordRequired };
    const previousCredential = previousCredentialRef.current;
    if (
      previousCredential.password !== nextCredential.password ||
      previousCredential.clearPassword !== nextCredential.clearPassword ||
      previousCredential.passwordRequired !== nextCredential.passwordRequired
    ) {
      previousCredentialRef.current = nextCredential;
      setCredentialGeneration((generation) => generation + 1);
    }
  }, [clearPassword, password, passwordRequired]);

  // StrictMode 會先 cleanup 再重跑 effect，必須在 setup 時重設為已掛載。
  useEffect(() => {
    mountedRef.current = true;
    return () => {
      mountedRef.current = false;
    };
  }, []);

  const previewStateKey = `${basePreviewStateKey}|credential-generation:${credentialGeneration}`;
  const operation = useSchemaOperation(workspaceKey, previewStateKey);
  const operationPreviewing = operation.phase === 'previewing' && operation.scopeSignature === previewStateKey;
  const schemaWritePending = operation.phase === 'creating' || operation.phase === 'unknown';
  const [schemaPhase, setSchemaPhase] = useState<LocalSchemaPhase>('idle');
  const [schemaStatements, setSchemaStatements] = useState<string[]>([]);
  const [schemaError, setSchemaError] = useState<string | null>(null);
  // A refusal is not a failure: the server asks for a confirmed preview first.
  const [schemaRefusal, setSchemaRefusal] = useState<SafeErrorMessage | null>(null);
  const [schemaExecuted, setSchemaExecuted] = useState(0);
  const schemaBusy = schemaPhase === 'previewing' || schemaPhase === 'creating' || operationPreviewing;
  const schemaControlsDisabled = schemaBusy || schemaActionsDisabled || schemaWritePending;

  useEffect(() => {
    setSchemaPhase('idle');
    setSchemaStatements([]);
    setSchemaError(null);
    setSchemaRefusal(null);
    setSchemaExecuted(0);
  }, [previewStateKey]);

  if (connector.kind === 'sqlite') {
    return null;
  }

  const handlePreviewSchema = async () => {
    if (schemaControlsDisabled) {
      return;
    }
    const token = beginSchemaOperation(workspaceKey, previewStateKey, 'previewing');
    if (!token) {
      return;
    }

    setSchemaPhase('previewing');
    setSchemaError(null);
    try {
      const result = await studioV2WorkspaceDatabaseAPI.generateSchema(true);
      const settled = settleSchemaOperation(token, 'idle');
      if (!settled || !mountedRef.current || !isSchemaOperationScopeCurrent(token)) {
        return;
      }
      setSchemaStatements(result.statements);
      setSchemaPhase('previewed');
    } catch (error) {
      const settled = settleSchemaOperation(token, 'idle');
      if (!settled || !mountedRef.current || !isSchemaOperationScopeCurrent(token)) {
        return;
      }
      setSchemaError(error instanceof Error ? error.message : String(error));
      setSchemaPhase('error');
    }
  };

  const handleCreateSchema = async () => {
    if (schemaActionsDisabled || schemaWritePending) {
      return;
    }
    const token = beginSchemaOperation(workspaceKey, previewStateKey, 'creating');
    if (!token) {
      return;
    }

    setSchemaPhase('creating');
    setSchemaError(null);
    setSchemaRefusal(null);
    try {
      const result = await studioV2WorkspaceDatabaseAPI.generateSchema(false);
      const settled = settleSchemaOperation(token, 'idle');
      if (!settled || !mountedRef.current || !isSchemaOperationScopeCurrent(token)) {
        return;
      }
      setSchemaStatements([]);
      setSchemaExecuted(result.executed);
      setSchemaPhase('done');
    } catch (error) {
      const httpFailure = isExplicitSchemaHTTPFailure(error);
      const settled = httpFailure
        ? settleSchemaOperation(token, 'idle')
        : settleSchemaOperation(token, 'unknown');
      if (!settled || !mountedRef.current || !isSchemaOperationScopeCurrent(token)) {
        return;
      }
      if (httpFailure) {
        const safe = getSafeErrorMessage(error, (key, options) => t(key, options));
        if (safe.code === 'SCHEMA_CONFIRMATION_REQUIRED') {
          setSchemaRefusal(safe);
          setSchemaPhase('idle');
          return;
        }
        setSchemaError(safe.message);
        setSchemaPhase('error');
      } else {
        setSchemaError(null);
        setSchemaPhase('unknown');
      }
    }
  };

  const operationPending = operation.phase === 'creating' || operation.phase === 'unknown';

  return (
    <div className="space-y-2 rounded-xl border border-gray-800/60 bg-gray-950/30 p-3">
      <div className="flex items-center justify-between gap-2">
        <div>
          <div className="text-xs text-gray-300">
            {t('step4.schema_section_title')}
          </div>
          <div className="mt-1 text-[11px] text-slate-500">
            {t('step4.schema_section_hint')}
          </div>
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            data-testid="schema-preview-btn"
            disabled={schemaControlsDisabled}
            onClick={handlePreviewSchema}
            className="rounded-md border border-gray-700 px-2 py-1 text-[11px] text-gray-300 transition-colors hover:bg-gray-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {schemaPhase === 'previewing' || operationPreviewing
              ? t('step4.schema_previewing')
              : t('step4.schema_preview_btn')}
          </button>
          <button
            type="button"
            data-testid="schema-create-btn"
            disabled={schemaControlsDisabled}
            onClick={handleCreateSchema}
            className="rounded-md border border-blue-700 px-2 py-1 text-[11px] text-blue-300 transition-colors hover:bg-blue-900/40 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {schemaPhase === 'creating' || operation.phase === 'creating'
              ? t('step4.schema_creating')
              : t('step4.schema_create_btn')}
          </button>
        </div>
      </div>

      {operationPending && (
        <p data-testid="schema-operation-pending" role="status" className="text-[11px] text-amber-300">
          {operation.phase === 'unknown'
            ? t('step4.schema_operation_unconfirmed')
            : t('step4.schema_operation_pending')}
        </p>
      )}

      {schemaPhase === 'previewed' && (
        schemaStatements.length === 0 ? (
          <p className="text-[11px] text-emerald-400">
            {t('step4.schema_no_changes')}
          </p>
        ) : (
          <pre className="max-h-32 overflow-auto whitespace-pre-wrap rounded-md bg-black/40 p-2 font-mono text-[10px] text-gray-300">
            {schemaStatements.join('\n')}
          </pre>
        )
      )}

      {schemaPhase === 'done' && (
        <p className="text-[11px] text-emerald-400">
          ✓ {t('step4.schema_done', { count: schemaExecuted })}
        </p>
      )}

      {schemaRefusal && (
        <div data-testid="schema-confirmation-required" role="status" className="break-words text-[11px] text-amber-300">
          <div>{t('step4.schema_confirmation_required')}</div>
          <div className="mt-1">{schemaRefusal.message}</div>
          {schemaRefusal.action && <div className="mt-1">{schemaRefusal.action}</div>}
        </div>
      )}

      {schemaPhase === 'error' && schemaError && (
        <p className="break-words text-[11px] text-rose-400">
          {t('step4.schema_error')}: {schemaError}
        </p>
      )}
    </div>
  );
}
