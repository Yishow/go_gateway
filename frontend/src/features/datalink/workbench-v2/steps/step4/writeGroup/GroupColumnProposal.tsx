import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';
import type { WriteGroupSqlDialect } from '../../../state/writeGroup/columns';
import { proposeColumns } from '../../../state/writeGroup/proposal';

export interface GroupColumnProposalProps {
  /** Tags the operator wants in the group that have no usable column. */
  unmatched: GroupCandidate[];
  existingColumns: string[];
  managed: boolean;
  /** Exact-value SQL dialect used when proposing a managed column type. */
  dialect?: WriteGroupSqlDialect;
}

/**
 * Offers the three honest ways out for Tags without a column — choose another
 * table, leave them out, or plan new columns — and shows the planned columns
 * labelled as a proposal. Nothing here is assigned or created.
 */
export const GroupColumnProposal: React.FC<GroupColumnProposalProps> = ({ unmatched, existingColumns, managed, dialect = 'portable' }) => {
  const { t } = useTranslation('workbench-v2');
  if (unmatched.length === 0) return null;
  const proposals = proposeColumns(unmatched, existingColumns, dialect);
  return (
    <section className="space-y-2 rounded border border-amber-800 p-3 text-xs text-slate-200" aria-label={t('step4.group.proposal.title')} data-testid="group-column-proposal">
      <h4 className="font-semibold text-amber-200">{t('step4.group.proposal.title', { count: unmatched.length })}</h4>
      <p>{managed ? t('step4.group.proposal.managed_intro') : t('step4.group.proposal.intro')}</p>
      <ul className="list-disc pl-5" data-testid="group-column-proposal-list">
        {proposals.map((proposal) => {
          const candidate = unmatched.find((entry) => entry.key === proposal.key);
          return (
            <li key={proposal.key}>
              {candidate?.label}: <span className="font-mono">{proposal.column} {proposal.sql_type}</span>{' '}
              <span className="text-amber-300">{t('step4.group.proposal.label')}</span>
            </li>
          );
        })}
      </ul>
      <p className="text-slate-400">{t('step4.group.proposal.options')}</p>
    </section>
  );
};
