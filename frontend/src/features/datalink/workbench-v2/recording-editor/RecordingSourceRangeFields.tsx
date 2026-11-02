import type { ReactNode } from 'react';
import { ParameterHelp } from '../components/ParameterHelp';
import type { FieldProps } from '../components/Field';
import './recordingSourceRangeFields.css';
/** Pure presentation shared by new and saved source forms. No draft or submission state. */
export function RecordingSourceRangeField({ label, hint, help, error, required, children, className = '' }: FieldProps) {
  return <div className={`recording-source-field ${className}`}>
    <label><span className="recording-source-field-label">{label}{required && <span aria-hidden="true"> *</span>}</span>{children}</label>
      {error ? <span role="alert" className="recording-source-field-hint">{error}</span> : hint && <span className="recording-source-field-hint">{hint}</span>}
    {help && <ParameterHelp compact help={help} />}
  </div>;
}
export function RecordingSourceRangeGrid({ children }: { children: ReactNode }) {
  return <div className="recording-source-range-grid">{children}</div>;
}
