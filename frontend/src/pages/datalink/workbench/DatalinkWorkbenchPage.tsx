import { useTranslation } from 'react-i18next';
import { WorkbenchActionDock } from './WorkbenchActionDock';
import { WorkbenchDeviceStep } from './WorkbenchDeviceStep';
import { WorkbenchHeaderBar } from './WorkbenchHeaderBar';
import { WorkbenchProvider, useWorkbench } from './WorkbenchProvider';
import { WorkbenchStepNavigator } from './WorkbenchStepNavigator';
import { LocalModbusBoard } from './LocalModbusBoard';
import { SourceCanvasSection } from './SourceCanvasSection';
import { TagBindingStudio } from './TagBindingStudio';
import { WORKBENCH_STEP_META } from './workbenchTypes';

function WorkbenchShell() {
  const { t } = useTranslation();
  const { activeStep } = useWorkbench();
  const activeStepMeta = WORKBENCH_STEP_META[activeStep];

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <WorkbenchHeaderBar />

      <div className="mx-auto grid w-full max-w-7xl gap-6 px-6 pb-8 xl:grid-cols-[240px_minmax(0,1fr)_320px]">
        <WorkbenchStepNavigator />

        <main className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 shadow-2xl shadow-slate-950/40">
          {activeStep === 'device' ? <WorkbenchDeviceStep /> : null}

          {activeStep === 'source' ? <SourceCanvasSection /> : null}

          {activeStep === 'tag' ? <TagBindingStudio /> : null}

          {activeStep === 'output' ? <LocalModbusBoard /> : null}

          {activeStep !== 'device' &&
          activeStep !== 'source' &&
          activeStep !== 'tag' &&
          activeStep !== 'output' ? (
            <section className="space-y-4 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
                {t('workbench.foundation.label')}
              </p>
              <div className="space-y-2">
                <h2 className="text-2xl font-semibold text-slate-50">
                  {t(activeStepMeta.labelKey)}
                </h2>
                <p className="max-w-2xl text-sm text-slate-300">
                  {t(activeStepMeta.placeholderKey)}
                </p>
              </div>
            </section>
          ) : null}
        </main>

        <WorkbenchActionDock />
      </div>
    </div>
  );
}

export default function DatalinkWorkbenchPage() {
  return (
    <WorkbenchProvider>
      <WorkbenchShell />
    </WorkbenchProvider>
  );
}
