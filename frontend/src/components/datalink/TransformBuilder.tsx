import { useState } from 'react';
import type { TransformStep, TransformType } from '../../types/datalink';

interface TransformBuilderProps {
  steps: TransformStep[];
  onChange: (steps: TransformStep[]) => void;
}

const TRANSFORM_TYPES: { type: TransformType; label: string; desc: string }[] = [
  { type: 'decode', label: 'Decode', desc: 'Binary decode (Endianness, etc.)' },
  { type: 'cast', label: 'Cast', desc: 'Type conversion (int to float, etc.)' },
  { type: 'scale', label: 'Scale', desc: 'Linear scaling (y = ax + b)' },
  { type: 'formula', label: 'Formula', desc: 'Mathematical expression' },
  { type: 'lookup', label: 'Lookup', desc: 'Value mapping dictionary' },
  { type: 'conditional', label: 'Conditional', desc: 'If/Else logic' },
];

export default function TransformBuilder({ steps, onChange }: TransformBuilderProps) {
  const [activeStepIndex, setActiveStepIndex] = useState<number | null>(null);

  const addStep = (type: TransformType) => {
    const newStep: TransformStep = {
      type,
      order: steps.length,
      params: getDefaultParams(type),
    };
    onChange([...steps, newStep]);
    setActiveStepIndex(steps.length);
  };

  const removeStep = (index: number) => {
    const newSteps = steps.filter((_, i) => i !== index).map((s, i) => ({ ...s, order: i }));
    onChange(newSteps);
    if (activeStepIndex === index) setActiveStepIndex(null);
  };

  const updateStepParams = (index: number, params: Record<string, any>) => {
    const newSteps = steps.map((s, i) => (i === index ? { ...s, params: { ...s.params, ...params } } : s));
    onChange(newSteps);
  };

  const getDefaultParams = (type: TransformType): Record<string, any> => {
    switch (type) {
      case 'scale': return { factor: 1, offset: 0 };
      case 'formula': return { expression: 'x' };
      case 'cast': return { target_type: 'float64' };
      case 'decode': return { byte_order: 'BigEndian', word_order: 'BigEndian' };
      case 'lookup': return { map: {} };
      default: return {};
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-2">
         <h3 className="text-sm font-semibold text-slate-300 uppercase tracking-wider">Pipeline Steps</h3>
         <span className="text-xs text-slate-500">{steps.length} Steps</span>
      </div>

      <div className="space-y-2 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
        {steps.map((step, index) => (
          <div key={index} className={`bg-slate-800 rounded-lg border transition-all duration-200 ${activeStepIndex === index ? 'border-blue-500 ring-1 ring-blue-500/50' : 'border-slate-700 hover:border-slate-600'}`}>
            <div 
                className="p-3 flex items-center justify-between cursor-pointer"
                onClick={() => setActiveStepIndex(activeStepIndex === index ? null : index)}
            >
              <div className="flex items-center space-x-3">
                <span className="w-6 h-6 rounded bg-slate-700 flex items-center justify-center text-xs font-mono text-slate-400">
                    {index + 1}
                </span>
                <span className="font-medium text-slate-200">
                    {TRANSFORM_TYPES.find(t => t.type === step.type)?.label || step.type}
                </span>
              </div>
              <button 
                  onClick={(e) => { e.stopPropagation(); removeStep(index); }}
                  className="p-1 hover:bg-red-500/10 hover:text-red-400 text-slate-500 rounded transition-colors"
                >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            {activeStepIndex === index && (
              <div className="p-3 border-t border-slate-700/50 bg-slate-900/30">
                 {renderStepConfig(step.type, step.params, (p) => updateStepParams(index, p))}
              </div>
            )}
          </div>
        ))}

        {steps.length === 0 && (
            <div className="text-center py-8 bg-slate-800/50 border border-slate-700 border-dashed rounded-lg">
                <p className="text-sm text-slate-500">No transform steps.</p>
                <p className="text-xs text-slate-600 mt-1">Add a step to process data.</p>
            </div>
        )}
      </div>

      <div className="grid grid-cols-3 gap-2 pt-2">
        {TRANSFORM_TYPES.map(t => (
            <button
                key={t.type}
                type="button"
                onClick={() => addStep(t.type)}
                className="px-2 py-2 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded-lg text-xs font-medium text-slate-300 transition-colors flex flex-col items-center gap-1 group"
                title={t.desc}
            >
                <span className="group-hover:text-blue-400 transition-colors">{t.label}</span>
            </button>
        ))}
      </div>
    </div>
  );
}

function renderStepConfig(type: TransformType, params: any, onChange: (p: any) => void) {
    switch (type) {
        case 'scale':
            return (
                <div className="grid grid-cols-2 gap-3">
                    <div>
                        <label className="text-xs text-slate-400 block mb-1">Factor (Mult)</label>
                        <input 
                            type="number" 
                            value={params.factor} 
                            onChange={e => onChange({ ...params, factor: parseFloat(e.target.value) })}
                            className="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-sm text-slate-200"
                        />
                    </div>
                    <div>
                        <label className="text-xs text-slate-400 block mb-1">Offset (Add)</label>
                        <input 
                            type="number" 
                            value={params.offset} 
                            onChange={e => onChange({ ...params, offset: parseFloat(e.target.value) })}
                            className="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-sm text-slate-200"
                        />
                    </div>
                </div>
            );
        case 'formula':
            return (
                <div>
                     <label className="text-xs text-slate-400 block mb-1">Expression (use 'x' as input)</label>
                     <input 
                        type="text" 
                        value={params.expression} 
                        onChange={e => onChange({ ...params, expression: e.target.value })}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-sm text-slate-200 font-mono"
                        placeholder="x * 10 + 5"
                    />
                </div>
            );
         case 'cast':
            return (
                <div>
                     <label className="text-xs text-slate-400 block mb-1">Target Type</label>
                     <select 
                        value={params.target_type} 
                        onChange={e => onChange({ ...params, target_type: e.target.value })}
                        className="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-sm text-slate-200"
                    >
                        <option value="int">Integer</option>
                        <option value="float">Float</option>
                        <option value="string">String</option>
                        <option value="bool">Boolean</option>
                    </select>
                </div>
            );
        // Implement other types as needed
        default:
            return <div className="text-xs text-slate-500 italic">No configuration available for this step type.</div>;
    }
}
