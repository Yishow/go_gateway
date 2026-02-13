import { useCallback, useEffect, useState } from 'react';
import { deviceAPI, pointAPI, tagAPI, mappingAPI } from '../../services/datalink';
import type {
  Device,
  Point,
  Tag,
  TransformStep,
  Mapping,
  MappingPreviewResponse,
  CreateMappingRequest,
  UpdateMappingRequest,
} from '../../types/datalink';
import TransformBuilder from './TransformBuilder';
import { useToast } from '../../contexts/ToastContext';
import { logger } from '../../utils/logger';

interface MappingCanvasProps {
  initialMapping?: Mapping | null;
  onSave: (data: CreateMappingRequest | UpdateMappingRequest) => Promise<void>;
  onCancel: () => void;
}

export default function MappingCanvas({ initialMapping, onSave, onCancel }: MappingCanvasProps) {
  // Data Sources
  const [devices, setDevices] = useState<Device[]>([]);
  const [points, setPoints] = useState<Point[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const { showError, showWarning } = useToast();

  // Selection State
  const [selectedDeviceId, setSelectedDeviceId] = useState<string>('');
  const [selectedPointId, setSelectedPointId] = useState<string>('');
  const [selectedTagId, setSelectedTagId] = useState<string>('');
  
  // Pipeline State
  const [steps, setSteps] = useState<TransformStep[]>([]);
  const [isEnabled, setIsEnabled] = useState(true);

  // UI State
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [previewResult, setPreviewResult] = useState<MappingPreviewResponse | null>(null);

  const loadPoints = useCallback(async (deviceId: string) => {
    try {
      const data = await pointAPI.list({ device_id: deviceId });
      setPoints(data);
    } catch (err) {
      logger.error("Failed to load points", err);
    }
  }, []);

  const loadInitialData = useCallback(async () => {
    setLoading(true);
    try {
      const [devicesData, tagsData] = await Promise.all([
        deviceAPI.list(),
        tagAPI.list({ status: 'active' }) // Only show active tags
      ]);
      setDevices(devicesData);
      setTags(tagsData);

      // If editing, we need to find the device of the current point to populate the point list
      if (initialMapping) {
          const point = await pointAPI.get(initialMapping.point_id);
          setSelectedDeviceId(point.device_id); // This triggers loadPoints
      }

    } catch (err) {
      logger.error("Failed to load data", err);
      showError("載入表單資料失敗");
    } finally {
      setLoading(false);
    }
  }, [initialMapping, showError]);

  useEffect(() => {
    loadInitialData();
  }, [loadInitialData]);

  // When device changes, fetch points
  useEffect(() => {
    if (selectedDeviceId) {
      loadPoints(selectedDeviceId);
    } else {
      setPoints([]);
    }
  }, [selectedDeviceId, loadPoints]);

  // Load initial mapping data if editing
  useEffect(() => {
    if (initialMapping && devices.length > 0 && tags.length > 0) {
      // Find device for the point (requires looking up point first, but for now we might need to rely on loaded data or load point details)
      // Note: In a real app we'd fetch the point details to know its device_id. 
      // Simplified here: We assume user selects device first. 
      // If editing, we just set point/tag IDs directly if available in lists.
      
      setSelectedPointId(initialMapping.point_id);
      setSelectedTagId(initialMapping.tag_id);
       
      let loadedSteps: TransformStep[] = [];
      try {
        loadedSteps = typeof initialMapping.transform_pipeline === 'string' 
            ? JSON.parse(initialMapping.transform_pipeline) 
            : initialMapping.transform_pipeline || [];
      } catch (e) {
        logger.error("Failed to parse pipeline", e);
      }
      setSteps(loadedSteps);
      setIsEnabled(initialMapping.enabled);
    }
  }, [initialMapping, devices, tags]); // Dependencies might need tuning in real scenario

  const handleSave = async () => {
    if (!selectedPointId || !selectedTagId) {
        showWarning("請選擇來源點位和目標標籤");
        return;
    }

    setSaving(true);
    try {
      await onSave({
        point_id: selectedPointId,
        tag_id: selectedTagId,
        transform_pipeline: steps,
        enabled: isEnabled,
      });
    } catch (error) {
       logger.error(error);
    } finally {
      setSaving(false);
    }
  };
  
  const handleTest = async () => {
       if (!selectedPointId) return;
       try {
           const res = await mappingAPI.preview({
               point_id: selectedPointId,
               transform_pipeline: steps
           });
           setPreviewResult(res);
        } catch(err: unknown) {
            const message =
              typeof err === 'object' && err !== null && 'message' in err
                ? String((err as { message?: string }).message ?? '未知錯誤')
                : '未知錯誤';
            showError(`預覽失敗: ${message}`);
        }
  };

  if (loading) return <div className="p-8 text-center text-slate-400">Loading editor resources...</div>;

  return (
    <div className="flex flex-col h-full bg-slate-900 text-slate-200">
      {/* Canvas Header */}
      <div className="p-6 border-b border-slate-800 flex justify-between items-center bg-slate-800/50 backdrop-blur">
         <div>
             <h2 className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-violet-400">
                 {initialMapping ? 'Edit Mapping Pipeline' : 'Design New Pipeline'}
             </h2>
             <p className="text-xs text-slate-500 mt-1">Connect Source Points to Destination Tags through transformation logic.</p>
         </div>
         <div className="flex items-center space-x-3">
             <div className="flex items-center space-x-2">
                 <button
                   type="button"
                   role="switch"
                   aria-checked={isEnabled}
                   aria-label="切換映射啟用狀態"
                   onClick={() => setIsEnabled(!isEnabled)}
                   className={`w-10 h-5 rounded-full p-1 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 ${isEnabled ? 'bg-emerald-500' : 'bg-slate-700'}`}
                 >
                     <div className={`w-3 h-3 bg-white rounded-full shadow-md transform transition-transform ${isEnabled ? 'translate-x-5' : 'translate-x-0'}`}></div>
                 </button>
                 <span className="text-sm font-medium text-slate-400">{isEnabled ? 'Active' : 'Paused'}</span>
             </div>
         </div>
      </div>

      {/* Main Canvas Area */}
      <div className="flex-1 overflow-auto p-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
          
          {/* LEFT: SOURCE */}
          <div className="space-y-4">
              <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700 h-full flex flex-col">
                  <div className="flex items-center space-x-2 mb-4 text-blue-400">
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                      </svg>
                      <h3 className="font-bold uppercase tracking-wider text-sm">Source (Input)</h3>
                  </div>

                  <div className="space-y-4 flex-1">
                      <div>
                          <label className="block text-xs font-medium text-slate-500 mb-1">Select Device</label>
                          <select 
                            value={selectedDeviceId}
                            onChange={(e) => setSelectedDeviceId(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm focus:ring-1 focus:ring-blue-500 outline-none transition-all"
                            disabled={!!initialMapping} // Lock device on edit for simplicity
                          >
                              <option value="">-- Choose Device --</option>
                              {devices.map(d => (
                                  <option key={d.id} value={d.id}>{d.name} ({d.protocol})</option>
                              ))}
                          </select>
                      </div>

                      <div>
                          <label className="block text-xs font-medium text-slate-500 mb-1">Select Point</label>
                          <select 
                             value={selectedPointId}
                             onChange={(e) => setSelectedPointId(e.target.value)}
                             className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm focus:ring-1 focus:ring-blue-500 outline-none transition-all disabled:opacity-50"
                             disabled={!selectedDeviceId || !!initialMapping}
                          >
                              <option value="">-- Choose Point --</option>
                              {points.map(p => (
                                  <option key={p.id} value={p.id}>{p.name} ({p.id.substring(0,8)})</option>
                              ))}
                          </select>
                      </div>

                      {selectedPointId && (
                          <div className="mt-4 p-3 bg-blue-500/5 rounded border border-blue-500/10">
                              <div className="text-xs text-blue-300">Selected Point Info</div>
                              {(() => {
                                  const p = points.find(point => point.id === selectedPointId);
                                  if (!p) return null;
                                  return (
                                      <div className="text-xs text-slate-400 mt-1 space-y-1">
                                          <div>Type: <span className="text-slate-200">{p.data_type}</span></div>
                                          <div>Address: <span className="text-slate-200">{p.address}</span></div>
                                      </div>
                                  );
                              })()}
                          </div>
                      )}
                  </div>
              </div>
          </div>

          {/* MIDDLE: PIPELINE */}
          <div className="space-y-4 relative">
             {/* Connector Lines (Visual Only) */}
             <div className="hidden lg:block absolute top-1/2 -left-6 w-6 h-0.5 bg-slate-700 border-t border-dashed border-slate-500"></div>
             <div className="hidden lg:block absolute top-1/2 -right-6 w-6 h-0.5 bg-slate-700 border-t border-dashed border-slate-500"></div>

              <div className="bg-slate-800/80 rounded-xl p-4 border border-slate-700 h-full flex flex-col shadow-xl shadow-black/20">
                  <div className="flex items-center space-x-2 mb-4 text-violet-400">
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
                      </svg>
                      <h3 className="font-bold uppercase tracking-wider text-sm">Processing Pipeline</h3>
                  </div>

                  <div className="flex-1 overflow-hidden">
                      <TransformBuilder steps={steps} onChange={setSteps} />
                  </div>
                  
                  {previewResult && (
                      <div className="mt-4 pt-4 border-t border-slate-700/50">
                          <div className="text-xs font-bold text-slate-300 mb-2">Last Test Result</div>
                          <div className="bg-slate-900 rounded p-2 text-xs font-mono space-y-1">
                              <div className="flex justify-between">
                                  <span className="text-slate-500">Input:</span>
                                  <span className="text-slate-200">{String(previewResult.raw_value)}</span>
                              </div>
                               <div className="flex justify-between">
                                  <span className="text-slate-500">Output:</span>
                                  <span className="text-emerald-400">{String(previewResult.final_value)}</span>
                              </div>
                              {previewResult.error && (
                                  <div className="text-red-400 mt-1 pt-1 border-t border-slate-800">
                                      Error: {previewResult.error}
                                  </div>
                              )}
                          </div>
                      </div>
                  )}

                  <button 
                      onClick={handleTest}
                      disabled={!selectedPointId}
                      className="mt-4 w-full py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 disabled:cursor-not-allowed text-xs font-bold tracking-wider rounded uppercase transition-colors"
                  >
                      Test Pipeline with Live Value
                  </button>
              </div>
          </div>

          {/* RIGHT: DESTINATION */}
          <div className="space-y-4">
              <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700 h-full flex flex-col">
                   <div className="flex items-center space-x-2 mb-4 text-emerald-400">
                      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                      </svg>
                      <h3 className="font-bold uppercase tracking-wider text-sm">Destination (Output)</h3>
                  </div>

                  <div className="space-y-4 flex-1">
                      <div>
                          <label className="block text-xs font-medium text-slate-500 mb-1">Select Tag</label>
                          <select 
                            value={selectedTagId}
                            onChange={(e) => setSelectedTagId(e.target.value)}
                            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm focus:ring-1 focus:ring-emerald-500 outline-none transition-all"
                            disabled={!!initialMapping}
                          >
                              <option value="">-- Choose Tag --</option>
                              {tags.map(t => (
                                  <option key={t.id} value={t.id}>{t.key}</option>
                              ))}
                          </select>
                      </div>

                      {selectedTagId && (
                          <div className="mt-4 p-3 bg-emerald-500/5 rounded border border-emerald-500/10">
                              <div className="text-xs text-emerald-300">Target Tag Info</div>
                               {(() => {
                                  const t = tags.find(tag => tag.id === selectedTagId);
                                  if (!t) return null;
                                  return (
                                      <div className="text-xs text-slate-400 mt-1 space-y-1">
                                          <div>Key: <span className="text-slate-200">{t.key}</span></div>
                                          <div>Name: <span className="text-slate-200">{t.display_name}</span></div>
                                          <div>Type: <span className="text-slate-200">{t.data_type}</span></div>
                                      </div>
                                  );
                              })()}
                          </div>
                      )}
                  </div>
              </div>
          </div>
      </div>

      {/* Footer Actions */}
      <div className="p-6 border-t border-slate-800 bg-slate-900 flex justify-end space-x-4">
          <button 
              onClick={onCancel}
              className="px-6 py-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
          >
              Cancel
          </button>
          <button 
              onClick={handleSave}
              disabled={saving}
              className="px-6 py-2 rounded-lg bg-blue-600 hover:bg-blue-500 text-white font-medium shadow-lg shadow-blue-500/20 disabled:opacity-50 disabled:cursor-wait transition-all flex items-center space-x-2"
          >
              {saving && <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>}
              <span>{initialMapping ? 'Update Mapping' : 'Create Mapping'}</span>
          </button>
      </div>
    </div>
  );
}
