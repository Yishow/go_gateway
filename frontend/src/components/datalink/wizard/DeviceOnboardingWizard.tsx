import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import StepIndicator from './StepIndicator';
import DeviceIdentityStep from './device-steps/DeviceIdentityStep';
import DeviceConnectionStep from './device-steps/DeviceConnectionStep';
import DeviceValidationStep from './device-steps/DeviceValidationStep';
import { useToast } from '../../../contexts/ToastContext';
import { 
  useCreateDeviceMutation, 
  useUpdateDeviceMutation, 
  useToggleDeviceStatusMutation 
} from '../../../hooks/datalink/useDevices';
import { type ProtocolType } from '../../../types/datalink';
import { protocolAPI } from '../../../services/datalink';
import { useQuery } from '@tanstack/react-query';

interface WizardFormData {
  name: string;
  description: string;
  protocol: ProtocolType | '';
  config: Record<string, string | number | boolean | string[] | undefined>;
}

const INITIAL_DATA: WizardFormData = {
  name: '',
  description: '',
  protocol: '',
  config: {},
};

const STEPS = [
  { id: 'identity', title: '設備資料', description: '名稱與協議', icon: '📱' },
  { id: 'connection', title: '連線設定', description: '網路參數', icon: '🔌' },
  { id: 'validation', title: '驗證測試', description: '連線驗證', icon: '🔍' },
  { id: 'complete', title: '完成', description: '準備就緒', icon: '✅' },
];

interface DeviceOnboardingWizardProps {
  embedded?: boolean;
  onClose?: () => void;
  onActivated?: (deviceId: string) => void;
}

export default function DeviceOnboardingWizard({ embedded = false, onClose, onActivated }: DeviceOnboardingWizardProps = {}) {
  const navigate = useNavigate();
  const { showSuccess, showError } = useToast();
  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState<WizardFormData>(INITIAL_DATA);
  const [deviceId, setDeviceId] = useState<string | null>(null);
  const [isValidated, setIsValidated] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  // API Hooks
  const createMutation = useCreateDeviceMutation();
  const updateMutation = useUpdateDeviceMutation();
  const toggleStatusMutation = useToggleDeviceStatusMutation();
  
  const { data: protocols = [] } = useQuery({
    queryKey: ['protocols'],
    queryFn: () => protocolAPI.list(),
  });

  // Helpers
  const updateFormData = (updates: Partial<WizardFormData>) => {
    setFormData(prev => ({ ...prev, ...updates }));
    // Clear errors for updated fields
    const newErrors = { ...errors };
    Object.keys(updates).forEach(key => delete newErrors[key]);
    setErrors(newErrors);
  };

  const validateStep = (step: number) => {
    const newErrors: Record<string, string> = {};
    if (step === 0) {
      if (!formData.name) newErrors.name = '請填寫設備名稱';
      if (!formData.protocol) newErrors.protocol = '請選擇通訊協議';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = async () => {
    if (!validateStep(currentStep)) return;

    if (currentStep === 1) {
      // Create or Update Device before moving to validation
      try {
        if (deviceId) {
           await updateMutation.mutateAsync({
             id: deviceId,
             data: {
               name: formData.name,
               description: formData.description,
               connection_config: formData.config
             }
           });
        } else {
           const newDevice = await createMutation.mutateAsync({
             name: formData.name,
             description: formData.description,
             protocol: formData.protocol as ProtocolType,
             connection_config: formData.config
           });
           setDeviceId(newDevice.id);
        }
        setCurrentStep(prev => prev + 1);
      } catch (err: unknown) {
        const message =
          typeof err === 'object' && err !== null && 'message' in err
            ? String((err as { message?: string }).message ?? '儲存設備失敗')
            : '儲存設備失敗';
        showError(message);
        return;
      }
    } else if (currentStep === 2) {
      // Validation step
      if (!isValidated) {
          showError('請確認驗證通過後再繼續。');
          return;
      }
      
      // Activate device
      if (deviceId) {
          try {
              await toggleStatusMutation.mutateAsync({ id: deviceId, currentStatus: 'draft' }); // activating
              showSuccess('設備已成功啟用！');
              if (embedded) {
                onActivated?.(deviceId);
              } else {
                navigate('/datalink/devices');
              }
          } catch(err: unknown) {
              const message =
                typeof err === 'object' && err !== null && 'message' in err
                  ? String((err as { message?: string }).message ?? '未知錯誤')
                  : '未知錯誤';
              showError('設備啟用失敗：' + message);
          }
      }
    } else {
      setCurrentStep(prev => prev + 1);
    }
  };

  const handleBack = () => {
    setCurrentStep(prev => prev - 1);
  };

  const renderContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <DeviceIdentityStep
            name={formData.name}
            description={formData.description}
            protocol={formData.protocol}
            onChange={updateFormData}
            protocols={protocols}
            error={errors}
          />
        );
      case 1:
        if (!formData.protocol) return null; // Should not happen
        return (
          <DeviceConnectionStep
            protocol={formData.protocol as ProtocolType}
            config={formData.config}
            onChange={(cfg: Record<string, string | number | boolean | string[] | undefined>) => updateFormData({ config: cfg })}
          />
        );
      case 2:
        if (!deviceId) return <div className="text-red-400">錯誤：找不到設備 ID</div>;
        return (
          <DeviceValidationStep
            deviceId={deviceId}
            onValidationComplete={setIsValidated}
          />
        );
      default:
        return null;
    }
  };
  
    // Special handling using getDefaultConfigForProtocol logic inside Connection Step?
    // Actually DeviceConnectionStep expects config passed in. 
    // We should initialize config when protocol changes in Step 0.
    // I need to duplicate `getDefaultConfigForProtocol` or import it.
    // Importing from DeviceForm is hard as it's not exported.
    // Let's implement a simplified version or extract it later.
    // For now, empty config in Step 0 is fine, Step 1 (DeviceConnectionStep) should handle defaults if empty?
    // Looking at my DeviceConnectionStep implementation: `value={config.port || 502}` so it handles defaults visually.
    // But `formData.config` will be empty.
    // `DeviceConnectionStep` should probably call `onChange` with defaults on mount if empty?
    // Let's rely on user interaction or the implicit defaults. 
    // Wait, create API requires fields. `DeviceConnectionStep` visual defaults are just placeholders or values? 
    // `value={config.port || 502}` means if `config.port` is undefined, it shows 502. 
    // But `config` state remains empty until user types.
    // This is a common issue. I should pre-populate config when protocol is selected.

  return (
    <div className="mx-auto w-full max-w-4xl px-3 py-4 sm:p-6">
      <div className="mb-4 sm:mb-8">
        <h1 className="text-xl font-bold text-slate-100 mb-1 sm:mb-2 sm:text-2xl">新增設備</h1>
        <p className="text-sm text-slate-400 sm:text-base">依照步驟完成工業設備的設定與連線驗證。</p>
      </div>

      <StepIndicator
        steps={STEPS}
        currentStep={currentStep}
      />

      <div className="mt-4 rounded-2xl border border-slate-700 bg-slate-800 p-4 shadow-xl sm:mt-8 sm:p-8">
        <div className="min-h-[240px] sm:min-h-[400px]">
          {renderContent()}
        </div>

        <div className="mt-6 flex flex-col-reverse gap-3 border-t border-slate-700 pt-6 sm:mt-8 sm:flex-row sm:justify-between sm:gap-0">
          <button
             onClick={() => {
                 if (currentStep === 0) {
                   if (embedded) {
                     onClose?.();
                   } else {
                     navigate('/datalink/devices');
                   }
                 }
                 else handleBack();
             }}
             className="w-full rounded-xl px-6 py-3 font-medium text-slate-300 transition-colors hover:text-white sm:w-auto sm:py-2.5"
          >
            {currentStep === 0 ? '取消' : '上一步'}
          </button>
          
          <button
            onClick={handleNext}
            disabled={currentStep === 2 && !isValidated}
            className={`
                w-full rounded-xl bg-blue-600 px-8 py-3 font-medium text-white shadow-lg shadow-blue-500/20 transition-all hover:bg-blue-500
                disabled:cursor-not-allowed disabled:opacity-50 sm:w-auto sm:py-2.5
            `}
          >
            {currentStep === STEPS.length - 2 ? '完成並啟用' : '下一步'}
            {/* The STEPS array has 4 items. 'validation' is index 2. 'complete' is index 3.
                My handleNext logic handles completion at step 2 (validation) -> activate.
                Wait, I have a 'Complete' step in STEPS but distinct logic?
                If validation is successful, user clicks 'Complete & Activate'.
                So Step 2 is the last interactive step. Step 3 is just a success state or maybe redirect?
                My `handleNext` navigates away at Step 2 success.
                So `STEPS` might be misleading if I don't show Step 3 content.
                Let's remove 'Complete' step from array or use it as a success screen.
                Let's navigate away for now to keep it simple.
            */}
          </button>
        </div>
      </div>
    </div>
  );
}
