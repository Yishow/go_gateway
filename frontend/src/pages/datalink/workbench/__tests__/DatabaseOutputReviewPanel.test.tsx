import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type {
  SourceRuleCandidateSetView,
  SourceRuleDatabaseOutputCandidateView,
} from '../../../../types/sourceRuleCandidates';
import { DatabaseOutputReviewPanel } from '../DatabaseOutputReviewPanel';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, string>) => {
      if (key === 'workbench.output.database.reviewSurface.groupedSuggestion') {
        return `列建議：${options?.groupKey} → ${options?.columnName}`;
      }
      if (key === 'workbench.output.database.reviewSurface.singleSuggestion') {
        return `單列建議：${options?.columnName}`;
      }
      return key;
    },
  }),
}));

const reviewSet: SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView> = {
  status: 'ready',
  candidates: [
    {
      id: 'grouped-candidate',
      identity: {
        source_rule_id: 'rule-1',
        candidate_type: 'database_outputs',
        candidate_kind: 'database_output',
        derived_from_rule_address: '40001',
      },
      proposed_signature: 'sig-1',
      address: '40001',
      point_id: 'point-1',
      tag_key: 'meter/A1',
      display_name: 'Flow A1',
      data_type: 'int16',
      status: 'ready',
      group_key: 'line',
      column_name: 'flow_kw',
    },
    {
      id: 'single-candidate',
      identity: {
        source_rule_id: 'rule-1',
        candidate_type: 'database_outputs',
        candidate_kind: 'database_output',
        derived_from_rule_address: '40002',
      },
      proposed_signature: 'sig-2',
      address: '40002',
      point_id: 'point-2',
      tag_key: 'meter/kw',
      display_name: 'Flow KW',
      data_type: 'int16',
      status: 'ready',
      group_key: null,
      column_name: 'kw',
    },
  ],
};

describe('DatabaseOutputReviewPanel', () => {
  it('shows grouped and single-row suggestions from the review set', () => {
    render(
      <DatabaseOutputReviewPanel
        reviewRuleId="rule-1"
        reviewRevisionId="rev-1"
        reviewSet={reviewSet}
        reviewLoading={false}
        selectedConnectorId=""
        selectedConnectorName={null}
        selectedTableKey=""
        selectedTagId=""
      />,
    );

    expect(screen.getByTestId('database-review-grouping-grouped-candidate')).toHaveTextContent(
      '列建議：line → flow_kw',
    );
    expect(screen.getByTestId('database-review-grouping-single-candidate')).toHaveTextContent(
      '單列建議：kw',
    );
  });
});
