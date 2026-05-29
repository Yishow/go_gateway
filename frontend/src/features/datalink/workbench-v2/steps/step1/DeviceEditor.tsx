import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Device, ProtocolId } from '../../state/types';
import { SectionCard } from '../../components/SectionCard';
import { Field } from '../../components/Field';
import { Input, Textarea } from '../../components/inputs';
import { ProtocolSelector } from './ProtocolSelector';
import { ConnectionConfigForm } from './ConnectionConfigForm';

export interface DeviceEditorProps {
  /** 當前正在編輯的設備資料 */
  device: Device;
  /** 更新設備基本屬性的 callback */
  onUpdate: (patch: Partial<Device>) => void;
  /** 更新設備配置參數的 callback */
  onUpdateConfig: (patch: Record<string, any>) => void;
  /** 變更設備協議的 callback */
  onChangeProtocol: (protocol: ProtocolId) => void;
}

/**
 * 設備參數編輯器元件 (DeviceEditor)
 * 
 * 落地需求 **Protocol switching with reset** / **Connection config form per protocol**。
 * 作為 Step 1 左側的主要編輯工作區，統合設備名稱、描述、通訊協議類型選擇，以及動態切換的連線參數表單。
 */
export const DeviceEditor: React.FC<DeviceEditorProps> = ({
  device,
  onUpdate,
  onUpdateConfig,
  onChangeProtocol,
}) => {
  const { t } = useTranslation('workbench-v2');

  return (
    <SectionCard
      title={t('step1.panels.editor')}
      icon="sliders"
      data-testid="device-editor"
    >
      <div className="space-y-5">
        {/* 設備名稱與描述 */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Field label={t('step1.fields.device_name')} required>
            <Input
              type="text"
              value={device.name}
              onChange={(e) => onUpdate({ name: e.target.value })}
              placeholder="e.g. PLC-01"
              data-testid="editor-input-name"
            />
          </Field>

          <Field label={t('step1.fields.device_desc')}>
            <Textarea
              value={device.description}
              onChange={(e) => onUpdate({ description: e.target.value })}
              placeholder={t('step1.hints.description_placeholder')}
              rows={1}
              className="resize-none"
              data-testid="editor-input-desc"
            />
          </Field>
        </div>

        {/* 協議選擇器 */}
        <div className="space-y-2 border-t border-slate-800/80 pt-4">
          <span className="text-xs font-semibold text-slate-400 block">
            {t('step1.fields.select_protocol')}
          </span>
          <ProtocolSelector value={device.protocol} onChange={onChangeProtocol} />
          <span className="text-[10px] text-slate-500 block leading-normal">
            {t('step1.hints.protocol_reset_warn')}
          </span>
        </div>

        {/* 協議專屬連線配置 */}
        <div className="space-y-3 border-t border-slate-800/80 pt-4">
          <span className="text-xs font-semibold text-slate-400 block">
            {t('step1.fields.connection_settings')}
          </span>
          <ConnectionConfigForm
            protocol={device.protocol}
            config={device.config}
            onChange={onUpdateConfig}
          />
        </div>
      </div>
    </SectionCard>
  );
};
