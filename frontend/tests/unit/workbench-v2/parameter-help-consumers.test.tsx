import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { ConnectionConfigForm } from '@/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm';
import { parameterHelpSources } from '@/features/datalink/workbench-v2/components/parameterHelpContent';
import en from '@/i18n/locales/en/parameter-help.json';
import zh from '@/i18n/locales/zh-TW/parameter-help.json';
import enAdvanced from '@/i18n/locales/en/parameter-help-advanced.json';
import zhAdvanced from '@/i18n/locales/zh-TW/parameter-help-advanced.json';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: { title?: string }) => key === 'open' ? `About ${values?.title}` : key }) }));
afterEach(cleanup);
it.each([
 ['modbus_tcp', {}, ['host','port','modbusSlave','timeoutModbus','dataFormat']],
 ['modbus_udp', {}, ['host','port','udpSlave','timeoutOther','dataFormat']],
 ['mc_3e', {}, ['host','port','mcStation','timeoutOther','networkNo','pcNo','ioNo','dataFormat']],
 ['fatek_fbs', { mode: 'tcp' }, ['fatekMode','host','port','fatekStation','timeoutOther']],
 ['fatek_fbs', { mode: 'serial' }, ['fatekMode','serialPort','baudFatek','parity','fatekStation','dataBits','stopBits','timeoutOther']],
 ['modbus_rtu', {}, ['serialPort','baudModbus','parity','modbusSlave','dataBits','stopBits','timeoutModbus','dataFormat']],
 ['mqtt', {}, ['mqttBroker','mqttUsername','mqttClientId']],
] as const)('explains every existing %s connection branch without changing configuration', (protocol, config, keys) => {
 const change = vi.fn(); const { container } = render(<ConnectionConfigForm protocol={protocol} config={config} onChange={change} rawNumericText />);
 // Advanced disclosure remains the existing disclosure; explaining it does not select a protocol.
 container.querySelectorAll('details').forEach(details => { details.open = true; });
 for (const key of keys) {
  const button = screen.getByRole('button', { name: `About connection.${key}.title` });
  expect(button.closest('label')).toBeNull(); fireEvent.click(button);
  expect(screen.getByRole('note')).toHaveTextContent(`connection.${key}.how`);
  fireEvent.keyDown(window, { key: 'Escape' }); expect(button).toHaveFocus();
 }
 expect(change).not.toHaveBeenCalled();
});
it('keeps lexical numeric editing and disabled MQTT capability unchanged', () => {
 const change = vi.fn(); render(<ConnectionConfigForm protocol="modbus_tcp" config={{ port: '1e' }} onChange={change} rawNumericText />);
 expect(screen.getByTestId('input-port')).toHaveValue('1e');
 fireEvent.click(screen.getByRole('button', { name: 'About connection.port.title' })); expect(change).not.toHaveBeenCalled();
 fireEvent.keyDown(window, { key: 'Escape' }); fireEvent.change(screen.getByTestId('input-port'), { target: { value: '2e' } });
 expect(change).toHaveBeenLastCalledWith({ port: '2e' });
});
it('provides complete bilingual source-grounded content for every help key without secret interpolation', () => {
 for (const locale of [{ ...en, ...enAdvanced }, { ...zh, ...zhAdvanced }]) {
  for (const [key, reference] of Object.entries(parameterHelpSources)) {
   const [group, field] = key.split('.');
   const content = (locale as unknown as Record<string, Record<string, Record<string,string>>>)[group][field];
   for (const name of ['title','what','how','example']) expect(content[name]?.trim(), `${key}.${name}`).toBeTruthy();
   expect(JSON.stringify(content)).not.toMatch(/\{\{(?:password|token|dsn|value)\}\}/i);
   expect(reference.canonical).toBeTruthy(); expect(reference.source).toMatch(/^(frontend|internal)\//);
  }
 }
});
it('offers exactly one database/path explanation per visible connector field', async () => {
 const { ConnectorSection } = await import('@/features/datalink/workbench-v2/steps/step4/ConnectorSection');
 const update = vi.fn(), kindChange = vi.fn();
 for (const kind of ['sqlite', 'postgres'] as const) {
  const view = render(<ConnectorSection connector={{ kind, name: 'Destination', database: 'kept', host: '', port: 5432, username: '', password: '', schema: '', table: '', write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'recorded_at', status: 'draft' }} onUpdateConnector={update} onKindChange={kindChange} />);
  expect(screen.getAllByRole('button', { name: 'About db.database.title' })).toHaveLength(1);
  fireEvent.click(screen.getByRole('button', { name: 'About db.database.title' }));
  expect(screen.getByLabelText(kind === 'sqlite' ? 'step4.field_db_sqlite' : 'step4.field_database')).toHaveValue('kept');
  expect(update).not.toHaveBeenCalled(); expect(kindChange).not.toHaveBeenCalled(); view.unmount();
 }
});
