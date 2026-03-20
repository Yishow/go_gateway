import type { ReactNode } from 'react';

interface TestPageShellProps {
  children: ReactNode;
}

export default function TestPageShell({ children }: TestPageShellProps) {
  return (
    <div className="min-h-screen bg-slate-950 py-4 sm:py-6">
      <main className="px-4 sm:px-6 xl:px-8 2xl:px-10">{children}</main>
    </div>
  );
}
