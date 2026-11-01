import * as React from 'react';

interface CollapsibleSupportCardProps {
  title: string;
  testId?: string;
  subtitle: string;
  collapsedLabel: string;
  expandedLabel: string;
  summary?: React.ReactNode;
  defaultOpen?: boolean;
  children: React.ReactNode;
}

export function CollapsibleSupportCard({
  title,
  testId,
  subtitle,
  collapsedLabel,
  expandedLabel,
  summary,
  defaultOpen = false,
  children,
}: CollapsibleSupportCardProps) {
  const [open, setOpen] = React.useState(defaultOpen);

  return (
    <section data-testid={testId} className="rounded-2xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-4">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 space-y-1">
          <h3 className="text-sm font-semibold text-[var(--gw-text)]">{title}</h3>
          <p className="text-xs text-[var(--gw-muted)]">{subtitle}</p>
          {summary ? (
            <div className="pt-2 text-xs text-[var(--gw-text)]">
              {summary}
            </div>
          ) : null}
        </div>

        <button
          data-testid={testId ? `${testId}-toggle` : undefined}
          type="button"
          aria-expanded={open}
          onClick={() => setOpen((current) => !current)}
          className="shrink-0 rounded-lg border border-[var(--gw-control)] px-3 py-2 text-xs font-medium text-[var(--gw-text)] transition-colors hover:bg-[var(--gw-surface)]"
        >
          {open ? expandedLabel : collapsedLabel}
        </button>
      </div>

      <div hidden={!open} className="mt-4 border-t border-[var(--gw-border)] pt-4">
          {children}
      </div>
    </section>
  );
}
