import * as React from 'react';

/** Reports one merged "an operation is active" signal to the page and clears it when the panel goes away. */
export function useReportOperationActive(active: boolean, report?: (active: boolean) => void): void {
  const latest = React.useRef(report);
  latest.current = report;
  React.useEffect(() => { latest.current?.(active); }, [active]);
  React.useEffect(() => () => latest.current?.(false), []);
}
