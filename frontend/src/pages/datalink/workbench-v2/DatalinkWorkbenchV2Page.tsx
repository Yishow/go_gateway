import { WorkbenchV2Shell } from '../../../features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useActivateStudioV2WorkspaceMutation } from '../../../hooks/datalink/useStudioV2WorkspaceActivation';
import { useStudioV2WorkspaceQuery } from '../../../hooks/datalink/useStudioV2Workspace';
import { useStudioV2AutosaveState } from './useStudioV2AutosaveState';
import '../../../features/datalink/workbench-v2/styles/workbench-v2.css';

interface DatalinkWorkbenchV2PageProps {
  navigateTo?: (target: string) => void;
}

interface WorkbenchV2BootstrapStateProps {
  status: 'loading' | 'error';
}

function WorkbenchV2BootstrapState({ status }: WorkbenchV2BootstrapStateProps) {
  const isError = status === 'error';

  return (
    <div
      data-workbench-v2="true"
      data-testid={isError ? 'workbench-v2-bootstrap-error' : 'workbench-v2-bootstrap-loading'}
      role={isError ? 'alert' : 'status'}
      className="wbv2-bootstrap-state bg-canvas min-h-screen font-sans text-slate-100"
    >
      <section
        data-testid="workbench-v2-bootstrap-shell"
        className="wbv2-bootstrap-shell"
        aria-label={isError ? 'Studio V2 工作區載入失敗' : 'Studio V2 工作區載入中'}
      >
        <div className="flex items-start justify-between gap-6">
          <div>
            <p className="font-mono text-[11px] uppercase tracking-[0.16em] text-cyan-200/70">
              Studio V2
            </p>
            <h1 className="mt-2 text-2xl font-semibold leading-tight text-slate-50">
              {isError ? '無法載入工作區' : '正在整理工作區狀態'}
            </h1>
            <p className="mt-2 max-w-[34rem] text-sm leading-6 text-slate-400">
              {isError
                ? '無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。'
                : '正在讀取工作區、設備、接入規則與資料庫目標，完成後會直接進入設定流程。'}
            </p>
          </div>
          <div className={isError ? 'wbv2-bootstrap-mark is-error' : 'wbv2-bootstrap-mark'}>
            {isError ? 'ERR' : 'SYNC'}
          </div>
        </div>
        <div className="mt-8 grid gap-3 sm:grid-cols-3" aria-hidden="true">
          <span className="wbv2-skeleton-line" />
          <span className="wbv2-skeleton-line" />
          <span className="wbv2-skeleton-line" />
        </div>
      </section>
    </div>
  );
}

/**
 * Datalink Workbench V2 頁面入口
 * 
 * 落地設計決策：「路由策略：/studio/v2 並存，/studio 保留為 fallback」
 * 限制 scope 與狀態綁定。
 */
export default function DatalinkWorkbenchV2Page({
  navigateTo,
}: DatalinkWorkbenchV2PageProps) {
  const workspaceQuery = useStudioV2WorkspaceQuery();
  const autosave = useStudioV2AutosaveState(workspaceQuery.isSuccess);
  const activationMutation = useActivateStudioV2WorkspaceMutation();

  if (
    workspaceQuery.isError ||
    (workspaceQuery.isSuccess && !workspaceQuery.data) ||
    autosave.devicesQuery.isError ||
    autosave.rulesQuery.isError ||
    autosave.mappingsQuery.isError ||
    autosave.databaseConfigQuery.isError ||
    autosave.databaseTargetsQuery.isError
  ) {
    return <WorkbenchV2BootstrapState status="error" />;
  }

  if (
    workspaceQuery.isLoading ||
    (workspaceQuery.isSuccess && !autosave.workspaceHydrated) ||
    autosave.devicesQuery.isLoading ||
    autosave.rulesQuery.isLoading ||
    autosave.mappingsQuery.isLoading ||
    autosave.databaseConfigQuery.isLoading ||
    autosave.databaseTargetsQuery.isLoading
  ) {
    return <WorkbenchV2BootstrapState status="loading" />;
  }

  const workspace = workspaceQuery.data!;

  return (
    <div
      data-workbench-v2="true"
      data-testid="workbench-v2-root"
      data-workspace-id={workspace.id}
      data-workspace-status={workspace.status}
      className="bg-canvas min-h-screen font-sans text-slate-100 antialiased"
    >
      {autosave.draftLossWarning && (
        <div
          data-testid="workbench-v2-draft-loss-warning"
          role="alert"
          className="mx-4 mt-4 rounded-xl border border-amber-500/30 bg-amber-500/10 px-4 py-3 text-sm text-amber-200 sm:mx-5 lg:mx-7"
        >
          {autosave.draftLossWarning}
        </div>
      )}
      <WorkbenchV2Shell
        state={autosave.state}
        actions={autosave.actions}
        navigateTo={navigateTo}
        activateWorkspace={() => activationMutation.mutateAsync()}
      />
    </div>
  );
}
