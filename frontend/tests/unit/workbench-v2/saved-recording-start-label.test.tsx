import { fireEvent, render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { GroupLifecycleBar } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({ useWriteGroupLifecycleMutation: () => ({ isPending: false, mutateAsync: vi.fn() }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({ useWriteGroupLifecycleQuery: () => ({ refetch: vi.fn() }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions', () => ({ lifecycleMatches: () => false, SavedRecordingActions: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog', () => ({ SavedRecordingRemoveDialog: () => null }));
it.each(['basic', 'mixed'] as const)('labels the actual %s review destination without changing the single existing action', kind => {
 const base = longProfileGroup({ id: 'record', circuitId: 'c', circuitName: 'Circuit', interval: 15, points: [0, 1] });
 const group = kind === 'mixed' ? { ...base, recording_profile: { kind: 'mixed' as const, result_shape: 'member-calculation-v1' as const } } : base;
 const review = vi.fn();
 render(<GroupLifecycleBar group={group} workspaceId="ws-1" workspaceRevision="W2" canApply readonly={false} onChanged={vi.fn()} onConflict={vi.fn()} onReviewScopedStart={review} />);
 expect(screen.getByTestId('group-apply')).toHaveTextContent('start.reviewSaved');
 expect(screen.getByTestId('group-apply-note')).toHaveTextContent('start.reviewSavedNote');
 fireEvent.click(screen.getByTestId('group-apply')); expect(review).toHaveBeenCalledExactlyOnceWith(group);
});
