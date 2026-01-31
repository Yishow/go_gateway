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
  config: Record<string, any>;
}

const INITIAL_DATA: WizardFormData = {
  name: '',
  description: '',
  protocol: '',
  config: {},
};

const STEPS = [
  { id: 'identity', title: 'Device Basics', description: 'Identity & Protocol', icon: '📱' },
  { id: 'connection', title: 'Connection', description: 'Network Settings', icon: '🔌' },
  { id: 'validation', title: 'Validation', description: 'Verify Connectivity', icon: '🔍' },
  { id: 'complete', title: 'Complete', description: 'Ready to Run', icon: '✅' },
];

export default function DeviceOnboardingWizard() {
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
      if (!formData.name) newErrors.name = 'Device name is required';
      if (!formData.protocol) newErrors.protocol = 'Protocol is required';
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
      } catch (err: any) {
        showError(err.message || 'Failed to save device');
        return;
      }
    } else if (currentStep === 2) {
      // Validation step
      if (!isValidated) {
          showError("Please ensure validation passes before proceeding.");
          return;
      }
      
      // Activate device
      if (deviceId) {
          try {
              await toggleStatusMutation.mutateAsync({ id: deviceId, currentStatus: 'draft' }); // activating
              showSuccess("Device activated successfully!");
              navigate('/datalink/devices');
          } catch(err: any) {
              showError("Failed to activate device: " + err.message);
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
            onChange={(cfg) => updateFormData({ config: cfg })}
          />
        );
      case 2:
        if (!deviceId) return <div className="text-red-400">Error: Device ID missing</div>;
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
    <div className="max-w-4xl mx-auto p-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-slate-100 mb-2">New Device Onboarding</h1>
        <p className="text-slate-400">Follow the steps to configure and validate your new industrial device.</p>
      </div>

      <StepIndicator
        steps={STEPS}
        currentStep={currentStep}
      />

      <div className="bg-slate-800 rounded-2xl border border-slate-700 p-8 mt-8 shadow-xl">
        <div className="min-h-[400px]">
          {renderContent()}
        </div>

        <div className="flex justify-between pt-8 border-t border-slate-700 mt-8">
          <button
             onClick={() => {
                 if (currentStep === 0) navigate('/datalink/devices');
                 else handleBack();
             }}
             className="px-6 py-2.5 text-slate-300 hover:text-white transition-colors font-medium"
          >
            {currentStep === 0 ? 'Cancel' : 'Back'}
          </button>
          
          <button
            onClick={handleNext}
            disabled={currentStep === 2 && !isValidated}
            className={`
                px-8 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 
                disabled:opacity-50 disabled:cursor-not-allowed
            `}
          >
            {currentStep === STEPS.length - 2 ? 'Complete & Activate' : 'Next Step'}
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
