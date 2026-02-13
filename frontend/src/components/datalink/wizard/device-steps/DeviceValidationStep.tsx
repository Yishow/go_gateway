import { useCallback, useEffect, useState } from 'react';
import { useCheckReadinessMutation } from '../../../../hooks/datalink/useDevices';

interface DeviceValidationStepProps {
  deviceId: string;
  onValidationComplete: (success: boolean) => void;
}

interface ReadinessResult {
  device_id: string;
  status: 'ready' | 'warning' | 'error';
  checks: Array<{
    name: string;
    pass: boolean;
    message: string;
  }>;
}

export default function DeviceValidationStep({ deviceId, onValidationComplete }: DeviceValidationStepProps) {
  const checkMutation = useCheckReadinessMutation();
  const [result, setResult] = useState<ReadinessResult | null>(null);

  const runCheck = useCallback(async () => {
    try {
      const res = await checkMutation.mutateAsync(deviceId);
      setResult(res);
      const isPass = res.status === 'ready' || res.status === 'warning';
      onValidationComplete(isPass);
    } catch {
      // Create a pseudo-result for error display
      setResult({
        device_id: deviceId,
        status: 'error',
        checks: [{ name: 'API Check', pass: false, message: 'Failed to communicate with server' }]
      });
      onValidationComplete(false);
    }
  }, [checkMutation, deviceId, onValidationComplete]);

  useEffect(() => {
    if (deviceId && !result) {
      runCheck();
    }
  }, [deviceId, result, runCheck]);

  if (checkMutation.isPending) {
    return (
      <div className="flex flex-col items-center justify-center py-12 space-y-4">
        <div className="w-12 h-12 border-4 border-blue-500/30 border-t-blue-500 rounded-full animate-spin"></div>
        <p className="text-slate-400">Verifying device configuration...</p>
      </div>
    );
  }

  if (!result) return null;

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
      <div>
        <h3 className="text-lg font-medium text-slate-200 mb-4">Connection Verification</h3>
        
        <div className={`p-4 rounded-xl border mb-6 flex items-center space-x-4 ${
          result.status === 'ready' ? 'bg-green-500/10 border-green-500/30' :
          result.status === 'warning' ? 'bg-yellow-500/10 border-yellow-500/30' :
          'bg-red-500/10 border-red-500/30'
        }`}>
           <div className={`p-2 rounded-full ${
               result.status === 'ready' ? 'bg-green-500 text-white' :
               result.status === 'warning' ? 'bg-yellow-500 text-white' :
               'bg-red-500 text-white'
           }`}>
               {result.status === 'ready' && <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>}
               {result.status === 'warning' && <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>}
               {result.status === 'error' && <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>}
           </div>
           <div>
               <h4 className="font-medium text-slate-200 uppercase tracking-wide text-sm">{result.status}</h4>
               <p className="text-slate-400 text-sm">
                   {result.status === 'ready' ? 'Device is reachable and configured correctly.' :
                    result.status === 'warning' ? 'Device is reachable but has configuration warnings.' :
                    'Device configuration failed validation checks.'}
               </p>
           </div>
        </div>

        <div className="space-y-3">
            {result.checks.map((check, idx) => (
                <div key={idx} className="flex items-center justify-between p-3 bg-slate-900 rounded-lg border border-slate-700">
                    <div className="flex items-center space-x-3">
                        {check.pass ? (
                             <svg className="w-5 h-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
                        ) : (
                             <svg className="w-5 h-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                        )}
                        <span className="text-slate-300 font-medium">{check.name}</span>
                    </div>
                    <span className="text-sm text-slate-500">{check.message}</span>
                </div>
            ))}
        </div>

        {result.status !== 'ready' && (
            <div className="mt-6 flex justify-end">
                <button 
                    onClick={runCheck}
                    className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg text-sm font-medium transition-colors"
                >
                    Retry Verification
                </button>
            </div>
        )}
      </div>
    </div>
  );
}
