import GlobalStyles from '@mui/material/GlobalStyles';
import { alpha } from '@mui/material/styles';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';

const sentry = tokens.archetype.sentry;
const variantRoot = '[data-variant="v2-mui"]';

const outputPanels = [
  `${variantRoot} [data-testid="output-primary-anchor"]`,
  `${variantRoot} [data-testid="modbus-secondary-panels"]`,
  `${variantRoot} [data-testid="database-selected-tag"]`,
  `${variantRoot} [data-testid="database-secondary-panels"]`,
  `${variantRoot} [data-testid="schema-snapshot"]`,
  `${variantRoot} [data-testid="dry-run-results"]`,
  `${variantRoot} [data-testid="database-dry-run-results"]`,
].join(', ');

const actionButtons = [
  `${variantRoot} [data-testid="output-modbus-dry-run"]`,
  `${variantRoot} [data-testid="output-modbus-sync"]`,
  `${variantRoot} [data-testid="output-database-toggle-connector"]`,
  `${variantRoot} [data-testid="output-database-refresh-validation"]`,
  `${variantRoot} [data-testid="output-database-save-connector"]`,
  `${variantRoot} [data-testid="output-database-generate-schema"]`,
].join(', ');

const targetButtons = `${variantRoot} [data-testid="output-tag-chips"] > button`;
const registerSlots = `${variantRoot} [data-testid^="register-slot-"]`;
const schemaColumns = `${variantRoot} [data-testid^="schema-column-"]`;
const targetSwitchers = [
  `${variantRoot} [data-testid="output-primary-anchor"] button[aria-pressed]`,
].join(', ');

/**
 * Sentry Incident Desk GlobalStyles for the **Output** step.
 *
 * Keeps the shared Local Modbus / Database workboard intact while re-skinning
 * the existing panels, slot buttons, and action controls under the v2 shell.
 */
export function MuiWorkbenchOutputStyles() {
  return (
    <GlobalStyles
      styles={{
        [outputPanels]: {
          fontFamily: tokens.typography.family.sentryUi,
          border: `1px solid ${sentry.border}`,
          borderRadius: '14px',
          background: `linear-gradient(180deg, ${alpha(sentry.elevated, 0.98)} 0%, ${alpha(
            sentry.panel,
            0.96,
          )} 100%)`,
          boxShadow: tokens.treatment.ambientPurple,
        },

        [actionButtons]: {
          borderRadius: '13px',
          border: '1px solid #584674',
          backgroundColor: sentry.accentMuted,
          color: tokens.text.primary,
          fontFamily: tokens.typography.family.sentryUi,
          fontSize: '13px',
          fontWeight: 700,
          letterSpacing: '0.2px',
          textTransform: 'uppercase',
          boxShadow: tokens.treatment.insetButton,
        },

        [targetButtons]: {
          borderRadius: '999px',
          border: `1px solid ${sentry.border}`,
          backgroundColor: alpha(sentry.panel, 0.9),
          color: tokens.text.secondary,
          fontFamily: tokens.typography.family.sentryUi,
          fontWeight: 700,
        },

        [targetSwitchers]: {
          borderRadius: '10px',
          borderColor: sentry.border,
          fontFamily: tokens.typography.family.sentryUi,
        },

        [registerSlots]: {
          borderRadius: '12px',
          boxShadow: 'rgba(0, 0, 0, 0.1) 0px 10px 15px -3px',
        },

        [schemaColumns]: {
          borderRadius: '12px',
          boxShadow: 'rgba(0, 0, 0, 0.08) 0px 8px 12px -4px',
        },

        [`${variantRoot} [data-testid="register-map-canvas"]`]: {
          borderRadius: '14px',
          outline: 'none',
        },

        [`${variantRoot} [data-testid="schema-snapshot"]`]: {
          outline: 'none',
        },

        [`${variantRoot} [data-testid="output-primary-anchor"] p[class*="tracking-"]`]: {
          color: sentry.highlight,
          fontFamily: tokens.typography.family.sentryUi,
        },
      }}
    />
  );
}
