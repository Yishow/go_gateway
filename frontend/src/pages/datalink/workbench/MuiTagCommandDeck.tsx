import Box from '@mui/material/Box';
import { MuiWorkbenchTagStyles } from './MuiWorkbenchTagStyles';
import { MuiTagIncidentDesk } from './MuiTagIncidentDesk';

export function MuiTagCommandDeck() {
  return (
    <Box sx={{ display: 'flex', minHeight: 0, flex: 1, flexDirection: 'column', gap: 1.5 }}>
      <MuiWorkbenchTagStyles />
      <MuiTagIncidentDesk />
    </Box>
  );
}
