import { useState, useEffect } from 'react';
import { settingsAPI } from '../../services/datalink';
import type { SystemSettings, UpdateSystemSettingsRequest } from '../../types/datalink';
import SettingsForm from '../../components/datalink/SettingsForm';

export default function SettingsPage() {
  const [settings, setSettings] = useState<SystemSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saveSuccess, setSaveSuccess] = useState(false);

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
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000); // Reset success message
    } catch (err: any) {
      alert(`Failed to update settings: ${err.message}`);
    }
  };

  if (loading) {
    return (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-slate-400">Loading settings...</p>
        </div>
    );
  }

  if (error) {
    return (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4 text-red-400">
          Error: {error}
        </div>
    );
  }

  return (
    <div className="space-y-6 max-w-4xl">
      <div>
        <h2 className="text-2xl font-bold text-slate-100">System Settings</h2>
        <p className="text-slate-400">Configure global data collection parameters</p>
      </div>

      {saveSuccess && (
         <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-lg p-3 text-emerald-400">
             Settings saved successfully!
         </div>
      )}

      {settings && (
        <div className="bg-slate-800 rounded-xl border border-slate-700 p-6">
            <SettingsForm settings={settings} onSubmit={handleUpdate} />
        </div>
      )}
    </div>
  );
}
