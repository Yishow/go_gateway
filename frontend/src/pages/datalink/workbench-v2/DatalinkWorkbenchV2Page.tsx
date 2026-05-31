import { WorkbenchV2Shell } from '../../../features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useActivateStudioV2WorkspaceMutation } from '../../../hooks/datalink/useStudioV2WorkspaceActivation';
import { useStudioV2WorkspaceQuery } from '../../../hooks/datalink/useStudioV2Workspace';
import { useStudioV2AutosaveState } from './useStudioV2AutosaveState';
import '../../../features/datalink/workbench-v2/styles/workbench-v2.css';

interface DatalinkWorkbenchV2PageProps {
  navigateTo?: (target: string) => void;
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
    workspaceQuery.isLoading ||
    (workspaceQuery.isSuccess && !autosave.workspaceHydrated) ||
    autosave.devicesQuery.isLoading ||
    autosave.rulesQuery.isLoading ||
    autosave.databaseConfigQuery.isLoading ||
    autosave.databaseTargetsQuery.isLoading
  ) {
    return (
      <div
        data-testid="workbench-v2-bootstrap-loading"
        className="flex min-h-screen items-center justify-center bg-canvas font-sans text-slate-100"
      >
        載入 Studio V2 工作區...
      </div>
    );
  }

  if (
    workspaceQuery.isError ||
    !workspaceQuery.data ||
    autosave.devicesQuery.isError ||
    autosave.rulesQuery.isError ||
    autosave.databaseConfigQuery.isError ||
    autosave.databaseTargetsQuery.isError
  ) {
    return (
      <div
        data-testid="workbench-v2-bootstrap-error"
        className="flex min-h-screen items-center justify-center bg-canvas px-6 text-center font-sans text-slate-100"
      >
        無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。
      </div>
    );
  }

  return (
    <div
      data-workbench-v2="true"
      data-testid="workbench-v2-root"
      data-workspace-id={workspaceQuery.data.id}
      data-workspace-status={workspaceQuery.data.status}
      className="bg-canvas min-h-screen font-sans text-slate-100 antialiased"
    >
      <WorkbenchV2Shell
        state={autosave.state}
        actions={autosave.actions}
        navigateTo={navigateTo}
        activateWorkspace={() => activationMutation.mutateAsync()}
      />
    </div>
  );
}
