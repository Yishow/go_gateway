import { useState, useEffect } from 'react';
import { settingsAPI } from '../../services/datalink';
import type { SystemSettings, UpdateSystemSettingsRequest } from '../../types/datalink';
import SettingsForm from '../../components/datalink/SettingsForm';
import { useToast } from '../../contexts/ToastContext';

export default function SettingsPage() {
  const [settings, setSettings] = useState<SystemSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { showSuccess, showError } = useToast();

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    setLoading(true);
    try {
      const data = await settingsAPI.get();
      setSettings(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch settings');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async (data: UpdateSystemSettingsRequest) => {
    try {
      const updated = await settingsAPI.update(data);
      setSettings(updated);
      showSuccess('設定已成功儲存');
    } catch (err: any) {
      showError(`儲存設定失敗: ${err.message}`);
    }
  };

  return (
    <div className="space-y-8 max-w-5xl">
      {/* Header with Glow */}
      <div className="flex justify-between items-center relative z-0">
          <div className="absolute -left-10 -top-10 w-32 h-32 bg-orange-600/20 blur-3xl pointer-events-none"></div>
          <div>
            <h2 className="text-3xl font-bold text-slate-100 tracking-tight">System Settings</h2>
            <p className="text-slate-400 mt-1">Configure global parameters for data collection and storage</p>
          </div>
      </div>

      {loading ? (
        <div className="bg-slate-800 rounded-xl border border-slate-700 p-8 animate-pulse space-y-6">
            <div className="h-6 bg-slate-700 rounded w-1/4"></div>
            <div className="space-y-4">
                <div className="h-10 bg-slate-700 rounded w-full"></div>
                <div className="h-10 bg-slate-700 rounded w-full"></div>
            </div>
            <div className="h-10 bg-slate-700 rounded w-32 ml-auto"></div>
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400 flex items-center space-x-3">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
           <span>{error}</span>
        </div>
      ) : settings && (
        <div className="relative">
            <div className="bg-slate-800 rounded-2xl border border-slate-700 p-8 shadow-sm">
                <SettingsForm settings={settings} onSubmit={handleUpdate} />
            </div>
        </div>
      )}
    </div>
  );
}
