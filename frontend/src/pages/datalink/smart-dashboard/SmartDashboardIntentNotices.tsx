import type { TFunction } from 'i18next';

interface SmartDashboardIntentNoticesProps {
  showLegacyNotice: boolean;
  legacyMigrationTitle: string;
  legacyMigrationDescription: string;
  legacyMigrationDismissLabel: string;
  dismissLegacyNotice: () => void;
  showSectionIntentNotice: boolean;
  sectionIntentLabel: string;
  dismissSectionIntentNotice: () => void;
  showModalIntentNotice: boolean;
  modalIntentLabel: string;
  closeWorkflowModal: () => void;
  t: TFunction;
}

export default function SmartDashboardIntentNotices({
  showLegacyNotice,
  legacyMigrationTitle,
  legacyMigrationDescription,
  legacyMigrationDismissLabel,
  dismissLegacyNotice,
  showSectionIntentNotice,
  sectionIntentLabel,
  dismissSectionIntentNotice,
  showModalIntentNotice,
  modalIntentLabel,
  closeWorkflowModal,
  t,
}: SmartDashboardIntentNoticesProps) {
  return (
    <>
      {showLegacyNotice && (
        <section className="mx-3 mt-3 rounded-2xl border border-amber-300/30 bg-amber-500/10 px-4 py-3 text-amber-100 shadow-lg shadow-amber-900/10 transition-all duration-300 sm:mx-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <p className="text-xs uppercase tracking-wider text-amber-200/80">{legacyMigrationTitle}</p>
              <p className="text-sm">{legacyMigrationDescription}</p>
            </div>
            <button
              type="button"
              onClick={dismissLegacyNotice}
              className="rounded-lg border border-amber-300/40 bg-amber-500/20 px-3 py-1.5 text-xs font-medium hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-400"
            >
              {legacyMigrationDismissLabel}
            </button>
          </div>
        </section>
      )}

      {showSectionIntentNotice && (
        <section className="mx-3 mt-3 rounded-2xl border border-blue-300/30 bg-blue-500/10 px-4 py-3 text-blue-100 shadow-lg shadow-blue-900/10 transition-all duration-300 sm:mx-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <p className="text-xs uppercase tracking-wider text-blue-200/80">{t('smartDashboard.intentNotices.sectionRedirectTitle')}</p>
              <p className="text-sm">
                {t('smartDashboard.intentNotices.sectionRedirectDesc', { label: sectionIntentLabel })}
              </p>
            </div>
            <button
              type="button"
              onClick={dismissSectionIntentNotice}
              className="rounded-lg border border-blue-300/40 bg-blue-500/20 px-3 py-1.5 text-xs font-medium hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-400"
            >
              {t('common.close')}
            </button>
          </div>
        </section>
      )}

      {showModalIntentNotice && (
        <section className="mx-3 mt-3 rounded-2xl border border-indigo-300/30 bg-indigo-500/10 px-4 py-3 text-indigo-100 shadow-lg shadow-indigo-900/10 transition-all duration-300 sm:mx-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <p className="text-xs uppercase tracking-wider text-indigo-200/80">{t('smartDashboard.intentNotices.modalRedirectTitle')}</p>
              <p className="text-sm">
                {t('smartDashboard.intentNotices.modalRedirectDesc', { label: modalIntentLabel })}
              </p>
            </div>
            <button
              type="button"
              onClick={closeWorkflowModal}
              className="rounded-lg border border-indigo-300/40 bg-indigo-500/20 px-3 py-1.5 text-xs font-medium hover:bg-indigo-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-400"
            >
              {t('common.close')}
            </button>
          </div>
        </section>
      )}
    </>
  );
}
