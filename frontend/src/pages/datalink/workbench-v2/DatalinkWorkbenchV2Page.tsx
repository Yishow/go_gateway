import { useWorkbenchV2State } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { WorkbenchV2Shell } from '../../../features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import '../../../features/datalink/workbench-v2/styles/workbench-v2.css';

/**
 * Datalink Workbench V2 頁面入口
 * 
 * 落地設計決策：「路由策略：/studio/v2 並存，/studio 保留為 fallback」
 * 限制 scope 與狀態綁定。
 */
export default function DatalinkWorkbenchV2Page() {
  const { state, ...actions } = useWorkbenchV2State();

  return (
    <div
      data-workbench-v2="true"
      data-testid="workbench-v2-root"
      className="bg-canvas min-h-screen font-sans text-slate-100 antialiased"
    >
      <WorkbenchV2Shell state={state} actions={{ state, ...actions }} />
    </div>
  );
}
