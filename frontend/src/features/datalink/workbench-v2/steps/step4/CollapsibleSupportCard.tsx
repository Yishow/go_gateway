import * as React from 'react';

interface CollapsibleSupportCardProps {
  title: string;
  subtitle: string;
  collapsedLabel: string;
  expandedLabel: string;
  summary?: React.ReactNode;
  defaultOpen?: boolean;
  children: React.ReactNode;
}

export function CollapsibleSupportCard({
  title,
  subtitle,
  collapsedLabel,
  expandedLabel,
  summary,
  defaultOpen = false,
  children,
}: CollapsibleSupportCardProps) {
  const [open, setOpen] = React.useState(defaultOpen);

  return (
    <section className="rounded-2xl border border-gray-800 bg-gray-900/10 p-4 backdrop-blur-sm">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 space-y-1">
          <h3 className="text-sm font-semibold text-white">{title}</h3>
          <p className="text-xs text-gray-500">{subtitle}</p>
          {summary ? (
            <div className="pt-2 text-xs text-gray-300">
              {summary}
            </div>
          ) : null}
        </div>

        <button
          type="button"
          aria-expanded={open}
          onClick={() => setOpen((current) => !current)}
          className="shrink-0 rounded-lg border border-gray-700 px-3 py-2 text-xs font-medium text-gray-200 transition-colors hover:bg-gray-800"
        >
          {open ? expandedLabel : collapsedLabel}
        </button>
      </div>

      {open ? (
        <div className="mt-4 border-t border-gray-800/70 pt-4">
          {children}
        </div>
      ) : null}
    </section>
  );
}
