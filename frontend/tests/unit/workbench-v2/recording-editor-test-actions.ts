import { fireEvent, screen } from '@testing-library/react';
export async function openNewRecordingDevice() {
  if (!screen.queryByLabelText('picker.search')) fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(await screen.findByRole('button', { name: 'picker.newSource' }));
  fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
  await screen.findByTestId('recording-source-name');
}
export function addRecordingTemperature(addLabel = 'picker.add:1') {
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
  fireEvent.click(screen.getByRole('button', { name: addLabel }));
}
