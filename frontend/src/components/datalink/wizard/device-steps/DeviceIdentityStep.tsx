import { type ProtocolType } from '../../../../types/datalink';

interface DeviceIdentityStepProps {
  name: string;
  description: string;
  protocol: ProtocolType | '';
  onChange: (updates: { name?: string; description?: string; protocol?: ProtocolType }) => void;
  protocols: { type: string; name: string; description: string }[];
  error?: { name?: string; protocol?: string };
}

export default function DeviceIdentityStep({
  name,
  description,
  protocol,
  onChange,
  protocols,
  error
}: DeviceIdentityStepProps) {
  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
      <div>
        <h3 className="text-lg font-medium text-slate-200 mb-4">設備基本資料</h3>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              設備名稱 <span className="text-red-400">*</span>
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => onChange({ name: e.target.value })}
              className={`w-full bg-slate-900 border ${error?.name ? 'border-red-500' : 'border-slate-700'} rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 transition-colors`}
              placeholder="例：產線 PLC-01"
              autoFocus
            />
            {error?.name && <p className="text-sm text-red-400 mt-1">{error.name}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">
              描述（選填）
            </label>
            <textarea
              value={description}
              onChange={(e) => onChange({ description: e.target.value })}
              className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 transition-colors resize-none"
              rows={3}
              placeholder="選填描述..."
            />
          </div>
        </div>
      </div>

      <div>
        <h3 className="text-lg font-medium text-slate-200 mb-4">通訊協議</h3>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {protocols.map((p) => (
            <div
              key={p.type}
              onClick={() => onChange({ protocol: p.type as ProtocolType })}
              className={`
                cursor-pointer rounded-xl p-4 border transition-all duration-200 relative overflow-hidden group
                ${protocol === p.type 
                  ? 'bg-blue-600/10 border-blue-500 ring-1 ring-blue-500' 
                  : 'bg-slate-800/50 border-slate-700 hover:border-slate-600 hover:bg-slate-800'}
              `}
            >
              <div className="flex items-start justify-between">
                <div>
                  <div className="font-medium text-slate-200 mb-1">{p.name}</div>
                  <div className="text-xs text-slate-400 line-clamp-2">{p.description}</div>
                </div>
                {protocol === p.type && (
                  <div className="text-blue-400">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
        {error?.protocol && <p className="text-sm text-red-400 mt-2">{error.protocol}</p>}
      </div>
    </div>
  );
}
