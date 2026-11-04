import { useState } from 'react';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { GroupIdentityFields } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields';
import { newDraft, type EditorDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

it('offers discovered schemas and their user and gateway tables without discarding an existing draft', () => {
  const catalog = { status: 'ready' as const, defaultSchema: 'public', schemas: [
    { name: 'public', tables: ['customer_readings', 'gw_group_42', 'gw_other_owner', ' spaced_table '] },
    { name: 'archive', tables: ['history'] },
    { name: 'empty_schema', tables: [] },
  ], blockedTables: [{ schema: 'public', name: 'gw_other_owner' }], retry: vi.fn() };
  function Editor() {
    const [draft, setDraft] = useState<EditorDraft>(() => ({ ...newDraft({ connector_id: 'C', connector_revision: 'R', table_schema: 'public', table_name: 'gw_group_42' }), storage_strategy: 'managed' }));
    const props = { catalog, section: 'destination' as const, draft, readonly: false, locked: false,
      onChange: (patch: Partial<typeof draft>) => setDraft(current => ({ ...current, ...patch })),
      onStorageStrategy: () => undefined };
    return <><GroupIdentityFields {...props} /><output data-testid="draft">{JSON.stringify(draft)}</output></>;
  }
  render(<Editor />);
  const schema = screen.getByTestId('group-schema'), table = screen.getByTestId('group-table');
  expect(schema.tagName).toBe('SELECT'); expect(table.tagName).toBe('SELECT');
  expect(within(schema).getByRole('option', { name: 'empty_schema' })).toBeInTheDocument();
  expect(within(table).getByRole('option', { name: 'customer_readings' })).toBeInTheDocument();
  expect(within(table).getByRole('option', { name: 'gw_group_42' })).toBeInTheDocument();
  expect(within(table).getByRole('option', { name: /gw_other_owner/ })).toBeDisabled();
  expect(within(table).getByRole('option', { name: /spaced_table/ })).toBeDisabled();
  expect(within(table).queryByRole('option', { name: 'history' })).not.toBeInTheDocument();
  fireEvent.change(schema, { target: { value: (within(schema).getByRole('option', { name: 'archive' }) as HTMLOptionElement).value } });
  expect(screen.getByTestId('draft')).toHaveTextContent('"table_name":"gw_group_42"');
  expect(within(table).getByRole('option', { name: 'history' })).toBeInTheDocument();
  fireEvent.change(table, { target: { value: (within(table).getByRole('option', { name: 'history' }) as HTMLOptionElement).value } });
  expect(screen.getByTestId('draft')).toHaveTextContent('"table_name":"history"');
});
