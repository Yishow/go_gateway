import Box from '@mui/material/Box';
import { MuiOutputIncidentDesk } from './MuiOutputIncidentDesk';
import { MuiWorkbenchOutputStyles } from './MuiWorkbenchOutputStyles';

export function MuiOutputCommandDeck() {
  return (
    <Box sx={{ display: 'flex', minHeight: 0, flex: 1, flexDirection: 'column', gap: 1.5 }}>
      <MuiWorkbenchOutputStyles />
      <MuiOutputIncidentDesk />
    </Box>
  );
}
