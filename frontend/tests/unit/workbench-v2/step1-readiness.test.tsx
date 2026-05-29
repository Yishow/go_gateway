import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ReadinessStages } from '../../../src/features/datalink/workbench-v2/steps/step1/ReadinessStages';
import { getStagesForProtocol } from '../../../src/features/datalink/workbench-v2/state/protocols';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      // 模擬翻譯行為，傳回 key 本身或特定對映字串
      if (key === 'step1.groups.connect') return '連線通道';
      if (key === 'step1.groups.probe') return '協議握手';
      return key;
    },
  }),
}));

describe('ReadinessStages Component', () => {
  it('Modbus TCP 的群組分組應包含連線通道與協議握手', () => {
    const stages = getStagesForProtocol('modbus_tcp');
    render(<ReadinessStages stages={stages} testStatus="draft" />);

    // 驗證群組標題是否正確呈現
    expect(screen.getByText('連線通道')).toBeInTheDocument();
    expect(screen.getByText('協議握手')).toBeInTheDocument();

    // 驗證 connect group 包含的 stages
    const groupConnect = screen.getByTestId('group-connect');
    expect(groupConnect).toHaveTextContent('step1.stages.resolve');
    expect(groupConnect).toHaveTextContent('step1.stages.connect');

    // 驗證 probe group 包含的 stages
    const groupProbe = screen.getByTestId('group-probe');
    expect(groupProbe).toHaveTextContent('step1.stages.probe');
  });

  it('Modbus RTU 的群組分組應正確包含連線通道（開啟序列埠、Handshake）與協議握手', () => {
    const stages = getStagesForProtocol('modbus_rtu');
    render(<ReadinessStages stages={stages} testStatus="draft" />);

    expect(screen.getByText('連線通道')).toBeInTheDocument();
    expect(screen.getByText('協議握手')).toBeInTheDocument();

    const groupConnect = screen.getByTestId('group-connect');
    expect(groupConnect).toHaveTextContent('step1.stages.open_port');
    expect(groupConnect).toHaveTextContent('step1.stages.handshake');

    const groupProbe = screen.getByTestId('group-probe');
    expect(groupProbe).toHaveTextContent('step1.stages.probe');
  });

  it('MQTT 的群組分組應正確包含連線通道與協議訂閱', () => {
    const stages = getStagesForProtocol('mqtt');
    render(<ReadinessStages stages={stages} testStatus="draft" />);

    expect(screen.getByText('連線通道')).toBeInTheDocument();
    expect(screen.getByText('協議握手')).toBeInTheDocument();

    const groupConnect = screen.getByTestId('group-connect');
    expect(groupConnect).toHaveTextContent('step1.stages.resolve');
    expect(groupConnect).toHaveTextContent('step1.stages.connect');

    const groupProbe = screen.getByTestId('group-probe');
    expect(groupProbe).toHaveTextContent('step1.stages.subscribe');
  });
});
