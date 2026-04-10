import type { ReactNode } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { muiExperimentTheme } from '../../../styles/muiExperimentTheme';

type MuiWorkbenchShellProps = {
  children: ReactNode;
};

/**
 * Wraps the workbench subtree with the MUI v2 experiment theme.
 *
 * Scoped to the workbench only — other routes (TestPage, Gateway) remain
 * unaffected. CssBaseline applies the dark background and reset inside
 * this subtree.
 */
export function MuiWorkbenchShell({ children }: MuiWorkbenchShellProps) {
  return (
    <ThemeProvider theme={muiExperimentTheme}>
      <CssBaseline enableColorScheme />
      {children}
    </ThemeProvider>
  );
}
